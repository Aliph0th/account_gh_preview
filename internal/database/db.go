package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
	mu   sync.RWMutex
}

func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.initSchema(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		count TEXT NOT NULL DEFAULT '0'
	);
	`
	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) GetTokens() (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var count string
	err := db.conn.QueryRow("SELECT COALESCE(count, '0') FROM tokens WHERE id = 1").Scan(&count)
	if err == sql.ErrNoRows {
		return "0", nil
	}
	return count, err
}

func (db *DB) SetTokens(count string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, err := db.conn.Exec("INSERT OR REPLACE INTO tokens (id, count) VALUES (1, ?)", count)
	return err
}

func (db *DB) Close() error {
	return db.conn.Close()
}
