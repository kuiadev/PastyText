package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/kuiadev/pastytext/data"
)

const testDbFile = "pastytext_test.db"

func TestIndexRoute(t *testing.T) {
	server, _ := setupTest(t)
	defer teardownTest(server)

	s := httptest.NewServer(server.Handler)
	defer s.Close()

	req := httptest.NewRequest(http.MethodGet, s.URL+"/index.html", nil)
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	resp := w.Result()

	// Also checking for 301 Moved Permanently due to the redirect from /index.html to /
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMovedPermanently {
		t.Errorf("Expected status code 200 or 301, got %v", resp.StatusCode)
	}
}

func TestIDRoute(t *testing.T) {
	server, pts := setupTest(t)
	defer teardownTest(server)

	s := httptest.NewServer(server.Handler)
	defer s.Close()

	req := httptest.NewRequest(http.MethodGet, s.URL+"/id", nil)
	w := httptest.NewRecorder()
	pts.idHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	msg := make(map[string]string)
	err = json.Unmarshal(data, &msg)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	if msg["friendly_name"] == "" {
		t.Errorf("Expected friendly_name to be non-empty")
	}

	if msg["ipaddress"] == "" {
		t.Errorf("Expected ipaddress to be non-empty")
	}
}

func TestWebsocketRoute(t *testing.T) {
	server, _ := setupTest(t)
	defer teardownTest(server)

	s := httptest.NewServer(server.Handler)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, _, err := websocket.Dial(ctx, s.URL+"/ws", &websocket.DialOptions{
		Subprotocols: []string{subprotocol}})
	defer c.Close(websocket.StatusNormalClosure, "closing connection")

	if err != nil {
		t.Errorf("Failed to dial websocket: %v", err)
	}

	_, _, err = c.Read(ctx)
	if err != nil {
		t.Errorf("Failed to read message: %v", err)
	}
}

func TestPasteSubmission(t *testing.T) {
	server, _ := setupTest(t)
	defer teardownTest(server)

	s := httptest.NewServer(server.Handler)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c, _, err := websocket.Dial(ctx, s.URL+"/ws", &websocket.DialOptions{
		Subprotocols: []string{subprotocol}})
	defer c.Close(websocket.StatusNormalClosure, "closing connection")

	if err != nil {
		t.Errorf("Failed to dial websocket: %v", err)
	}

	readMsgChan := make(chan chanData)
	defer close(readMsgChan)

	msg := map[string]string{"user": "thorough-tester", "action": "add", "text": "hello world!"}
	err = wsjson.Write(ctx, c, msg)
	if err != nil {
		t.Errorf("Failed to write message: %v", err)
	}

	var pasteResp data.Paste = data.Paste{}
	// The first iteration would be empty because the db contained no data, so we need to read twice
	for i := range 2 {
		go func(readMsgChan chan<- chanData) {
			var message []data.Paste
			err = wsjson.Read(ctx, c, &message)

			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					err = ctx.Err()
				}
				readMsgChan <- chanData{content: "", err: err}
				return
			}
			readMsgChan <- chanData{content: message, err: nil}

		}(readMsgChan)

		chanResult := <-readMsgChan
		if chanResult.err != nil {
			t.Errorf("Failed to read message: %v", chanResult.err)
			if chanResult.err != context.DeadlineExceeded {
				t.Errorf("Deadline exceeded: %v", chanResult.err)
			}
		}
		if i == 1 {
			pr := chanResult.content.([]data.Paste)
			if len(pr) > 0 {
				pasteResp = pr[0]
			}
		}
	}

	if pasteResp.Content != "hello world!" {
		t.Errorf("Expected pasteResp[0].Content to be 'hello world!', got %v", pasteResp.Content)
	}
}

func TestPasteDeletion(t *testing.T) {
	server, _ := setupTest(t)
	defer teardownTest(server)

	s := httptest.NewServer(server.Handler)
	defer s.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c, _, err := websocket.Dial(ctx, s.URL+"/ws", &websocket.DialOptions{
		Subprotocols: []string{subprotocol}})
	defer c.Close(websocket.StatusNormalClosure, "closing connection")

	if err != nil {
		t.Errorf("Failed to dial websocket: %v", err)
	}

	readMsgChan := make(chan chanData)
	defer close(readMsgChan)

	msg := map[string]string{"user": "thorough-tester", "action": "add", "text": "hello world!"}
	err = wsjson.Write(ctx, c, msg)
	if err != nil {
		t.Errorf("Failed to write message: %v", err)
	}

	delmsg := map[string]interface{}{"id": 1, "action": "delete"}
	err = wsjson.Write(ctx, c, delmsg)
	if err != nil {
		t.Errorf("Failed to write message: %v", err)
	}

	var pasteResp data.Paste = data.Paste{}
	// The first iteration would be empty because the db contained no data, so we need to read twice
	for i := range 3 {
		go func(readMsgChan chan<- chanData) {
			var message []data.Paste
			err = wsjson.Read(ctx, c, &message)

			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					err = ctx.Err()
				}
				readMsgChan <- chanData{content: "", err: err}
				return
			}
			readMsgChan <- chanData{content: message, err: nil}

		}(readMsgChan)

		chanResult := <-readMsgChan
		if chanResult.err != nil {
			t.Errorf("Failed to read message: %v", chanResult.err)
			if chanResult.err != context.DeadlineExceeded {
				t.Errorf("Deadline exceeded: %v", chanResult.err)
			}
		}

		if i == 2 {
			pr := chanResult.content.([]data.Paste)

			if len(pr) > 0 {
				pasteResp = pr[0]
			}
		}
	}

	if pasteResp.Content != "" {
		t.Errorf("Expected pasteResp[0].Content to be empty, got %v", pasteResp.Content)
	}
}

func setupTest(t *testing.T) (*http.Server, *ptServer) {
	// Use a test database file
	os.Setenv("DB_FILE", testDbFile)

	pts, err := NewPtServer()
	if err != nil {
		t.Errorf("Failed to create server: %v", err)
		return nil, nil
	}

	server := &http.Server{
		Handler:      pts,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	return server, pts
}

func teardownTest(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Shutdown(ctx)

	os.Remove(testDbFile)
	os.Setenv("DB_FILE", "")
}
