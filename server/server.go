package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/kuiadev/pastytext/data"
	"github.com/mileusna/useragent"
)

// The subprotocol is a string that identifies the protocol that the server and client will use to communicate.
const subprotocol = "pastytextProtocol"

// The clientTimeout is the time that the server will wait for a message from the client.
const clientTimeout = time.Minute * 5

// ptServer is a struct that implements the http.Handler interface.
type ptServer struct {
	clients  map[*client]struct{}
	dbm      *data.Manager
	serveMux http.ServeMux
}

type client struct {
	message clientMessage
	conn    *websocket.Conn
	network string
	device  string
}

type clientMessage struct {
	Id      int    `json:"id"`
	User    string `json:"user"`
	Action  string `json:"action"`
	Text    string `json:"text"`
	Network string `json:"network"`
	Device  string `json:"device"`
}

type chanData struct {
	content interface{}
	err     error
}

func NewPtServer() (*ptServer, error) {
	dbm, err := data.NewManager()
	if err != nil {
		log.Fatalf("Failed to create data manager: %v\n", err)
		return nil, err
	}
	pt := &ptServer{
		clients: make(map[*client]struct{}),
		dbm:     dbm,
	}

	pt.serveMux.Handle("/", http.FileServer(http.Dir("./web")))
	pt.serveMux.HandleFunc("/id", pt.idHandler)
	pt.serveMux.HandleFunc("/ws", pt.joinHandler)

	return pt, nil
}

func (p *ptServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.serveMux.ServeHTTP(w, r)
}

func (p *ptServer) getRequestIP(r *http.Request) string {
	requestIP := r.Header.Get("X-forwarded-for")
	if requestIP == "" {
		requestIP = r.RemoteAddr
	}
	return strings.Split(requestIP, ":")[0]
}

func (p *ptServer) idHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/josn")
	idn := struct {
		Friendly_name string `json:"friendly_name"`
		IPaddress     string `json:"ipaddress"`
	}{Friendly_name: data.GenerateName(), IPaddress: p.getRequestIP(r)}

	idJson, err := json.Marshal(idn)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(idJson)
}

func (p *ptServer) joinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{subprotocol},
	})
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	defer conn.CloseNow()

	if conn.Subprotocol() != subprotocol {
		conn.Close(websocket.StatusPolicyViolation, fmt.Sprintf("Expected subprotocol %s, but got %s\n", subprotocol, conn.Subprotocol()))
		return
	}

	ua := useragent.Parse(r.UserAgent())
	c := &client{
		conn:    conn,
		message: clientMessage{},
		network: p.getRequestIP(r),
		device:  fmt.Sprintf("%s-%s", ua.OS, ua.Name),
	}

	p.clients[c] = struct{}{}
	p.joinClient(c)
}

// joinClient is a method that will be called when a new client connects to the server.
// c is a pointer to a client struct.
func (p *ptServer) joinClient(c *client) {
	sendErrMsgChan := make(chan error)
	defer close(sendErrMsgChan)

	//Send initial message to client
	go func() {
		pastes, err := p.dbm.GetPastes(c.network)
		if err != nil {
			sendErrMsgChan <- err
		} else {
			err := c.sendMessageToClient(pastes)
			sendErrMsgChan <- err
		}

	}()

	emsg := <-sendErrMsgChan
	if emsg != nil {
		log.Printf("error sending initial message to client: %v\n", emsg)
		c.conn.CloseNow()
		delete(p.clients, c)
	}

	//Read messages from client
	for {
		readMsgChan := make(chan chanData)
		defer close(readMsgChan)
		go c.readMessageFromClient(readMsgChan)

		var newClientMessage = clientMessage{}
		chanResult := <-readMsgChan
		if chanResult.err != nil {
			if chanResult.err != context.DeadlineExceeded {
				log.Printf("error reading message from client: %v\n", chanResult.err)
				c.conn.CloseNow()
				delete(p.clients, c)
				return
			} else {
				continue
			}
		}
		newClientMessage = chanResult.content.(clientMessage)

		if newClientMessage.Action == "add" {
			newClientMessage.Network = c.network
			newClientMessage.Device = c.device
			p.persistMessageFromClient(newClientMessage)
		} else {
			//The only other action is delete
			p.deletePaste(int64(newClientMessage.Id))
		}

		pastes, err := p.dbm.GetPastes(c.network)
		if err == nil {
			p.publishMessageToClients(pastes)
		}
	}
}

// persistMessageFromClient is a method that saves the message sent by the client to the database.
func (p *ptServer) persistMessageFromClient(msg clientMessage) {
	paste := data.Paste{
		User:      msg.User,
		Device:    msg.Device,
		Network:   msg.Network,
		Content:   msg.Text,
		CreatedAt: time.Now(),
	}
	_, err := p.dbm.InsertPaste(paste)
	if err != nil {
		log.Printf("error inserting paste: %v\n", err)
	}
}

func (p *ptServer) deletePaste(id int64) {
	err := p.dbm.DeletePaste(id)
	if err != nil {
		log.Printf("error deleting paste: %v\n", err)
	}
}

// publishMessageToClients is a method that sends a message to all clients.
func (p *ptServer) publishMessageToClients(pastes []data.Paste) {
	var wg sync.WaitGroup

	for c := range p.clients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.sendMessageToClient(pastes)
		}()
	}

	wg.Wait()
}

// readMessageFromClient is a method that reads messages from the client.
// msgChan is a channel that will receive the message.
func (c *client) readMessageFromClient(msgChan chan<- chanData) {
	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	var message clientMessage

	err := wsjson.Read(ctx, c.conn, &message)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			err = ctx.Err()
		}
		msgChan <- chanData{content: "", err: err}
		return
	}

	msgChan <- chanData{content: message, err: nil}
}

// sendMessageToClient is a method that sends a message to a client.
func (c *client) sendMessageToClient(pastes []data.Paste) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := wsjson.Write(ctx, c.conn, pastes)
	if err != nil {
		log.Printf("error writing message: %v\n", err)
		return err
	}
	return nil
}
