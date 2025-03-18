package data

import (
	"os"
	"testing"
	"time"
)

const testDbFile = "pastytext_test.db"

func TestNewManage(t *testing.T) {
	setupTest()
	defer teardownTest()

	manager, err := NewManager()
	if err != nil {
		t.Errorf("Failed to create new Manager: %v", err)
	}

	if manager.db == nil {
		t.Errorf("sql.DB of manager object nil")
	}

	if manager.db.Stats().OpenConnections == 0 {
		t.Errorf("Expected open connections to be 1 got: %v", manager.db.Stats().OpenConnections)
	}
}

func TestClose(t *testing.T) {
	setupTest()
	defer teardownTest()

	manager, err := NewManager()
	if err != nil {
		t.Errorf("Failed to create new Manager: %v", err)
	}

	manager.db.Close()

	if manager.db.Stats().OpenConnections > 0 {
		t.Errorf("Expected open connections to be 0 got: %v", manager.db.Stats().OpenConnections)
	}
}

func TestInsertPaste(t *testing.T) {
	setupTest()
	defer teardownTest()

	manager, err := NewManager()
	defer manager.db.Close()
	if err != nil {
		t.Errorf("Failed to create new Manager: %v", err)
	}

	paste := Paste{
		User:      "test User",
		Device:    "test-device",
		Network:   "test-network",
		Content:   "TestInsertPaste",
		CreatedAt: time.Now(),
	}

	i, err := manager.InsertPaste(paste)
	if err != nil {
		t.Errorf("Failed to insert new paste: %v", err)
	}

	// More than likely an error would happen before it reached this logic
	if i < 1 {
		t.Errorf("Insert id is invalid, the paste probably failed, id: %v", i)
	}
}

func TestGetPastes(t *testing.T) {
	setupTest()
	defer teardownTest()

	manager, err := NewManager()
	defer manager.db.Close()
	if err != nil {
		t.Errorf("Failed to create new Manager: %v", err)
	}

	paste := Paste{
		User:      "test User",
		Device:    "test-device",
		Network:   "test-network",
		Content:   "TestGetPastes",
		CreatedAt: time.Now(),
	}

	_, err = manager.InsertPaste(paste)
	if err != nil {
		t.Errorf("Failed to insert new paste: %v", err)
	}

	pastes, err := manager.GetPastes("test-network")
	if err != nil {
		t.Errorf("Failed to fetch pastes: %v", err)
	}

	if len(pastes) == 0 {
		t.Error("Expected pastes to not be empty")
	}

	if pastes[0].Content != "TestGetPastes" {
		t.Errorf("Expected paste content to 'TestGetPastes' got %v", pastes[0].Content)
	}
}

func setupTest() {
	os.Setenv("DB_FILE", testDbFile)
}

func teardownTest() {
	os.Remove(testDbFile)
	os.Setenv("DB_FILE", "")
}
