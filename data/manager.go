package data

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// create is a SQL query that creates the pastes table.
const create = `CREATE TABLE IF NOT EXISTS pastes (
	id INTEGER NOT NULL PRIMARY KEY,
	created_at DATETIME NOT NULL,
	user TEXT,
	device TEXT,
	content TEXT
);`

const dbfile string = "/var/pastytext/data/pastytext.db"

type Manager struct {
	db *sql.DB
}

// Paste is a struct that represents a paste.
type Paste struct {
	Id        int64
	CreatedAt time.Time
	User      string
	Device    string
	Content   string
}

func NewManager() (*Manager, error) {
	db, err := sql.Open("sqlite3", dbfile)
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
	res, err := m.db.Exec("INSERT INTO pastes (created_at, user, device, content) VALUES (?, ?, ?, ?)", p.CreatedAt, p.User, p.Device, p.Content)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// GetPastes returns all pastes from the database.
func (m *Manager) GetPastes() ([]Paste, error) {
	rows, err := m.db.Query("SELECT id, created_at, user, device, content FROM pastes ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []Paste
	for rows.Next() {
		var p Paste
		if err := rows.Scan(&p.Id, &p.CreatedAt, &p.User, &p.Device, &p.Content); err != nil {
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
