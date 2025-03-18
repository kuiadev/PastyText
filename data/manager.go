package data

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// create is a SQL query that creates the pastes table.
const create = `CREATE TABLE IF NOT EXISTS pastes (
	id INTEGER NOT NULL PRIMARY KEY,
	created_at DATETIME NOT NULL,
	network TEXT,
	user TEXT,
	device TEXT,
	content TEXT
);`

const defaultDbFile string = "pastytext.db"

type Manager struct {
	db *sql.DB
}

// Paste is a struct that represents a paste.
type Paste struct {
	Id        int64
	CreatedAt time.Time
	Network   string
	User      string
	Device    string
	Content   string
}

func NewManager() (*Manager, error) {
	dbFile := os.Getenv("DB_FILE")
	if dbFile == "" {
		dbFile = defaultDbFile
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(create); err != nil {
		return nil, err
	}

	return &Manager{db: db}, nil
}

func (m *Manager) Close() error {
	return m.db.Close()
}

// InsertPaste inserts a paste into the database.
func (m *Manager) InsertPaste(p Paste) (int64, error) {
	res, err := m.db.Exec("INSERT INTO pastes (created_at, network, user, device, content) VALUES (?, ?, ?, ?, ?)", p.CreatedAt, p.Network, p.User, p.Device, p.Content)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// GetPastes returns all pastes from the database.
func (m *Manager) GetPastes(network string) ([]Paste, error) {
	rows, err := m.db.Query("SELECT id, created_at, network, user, device, content FROM pastes WHERE network = ? ORDER BY created_at DESC", network)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var p Paste
		if err := rows.Scan(&p.Id, &p.CreatedAt, &p.Network, &p.User, &p.Device, &p.Content); err != nil {
			return nil, err
		}
		pastes = append(pastes, p)
	}

	return pastes, nil
}

// DeletePaste deletes a paste from the database based on its ID.
func (m *Manager) DeletePaste(id int64) error {
	_, err := m.db.Exec("DELETE FROM pastes WHERE id = ?", id)
	return err
}
