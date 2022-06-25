package db

import (
	"context"
	"fmt"

	"github.com/dayvillefire/pocsag-monitor/obj"
	"github.com/genjidb/genji"
)

var (
	DB *genji.DB

	opened bool
)

func OpenDB(fn string) (*genji.DB, error) {
	if opened {
		return &genji.DB{}, fmt.Errorf("database already open")
	}
	DB, err := genji.Open(fn)
	if err != nil {
		return &genji.DB{}, err
	}
	DB = DB.WithContext(context.Background())
	return DB, nil
}

func InitDB(db *genji.DB) error {
	err := db.Exec(`
    CREATE TABLE entry (
        timestamp       INT     PRIMARY KEY,
        cap             TEXT    NOT NULL,
        message         TEXT    NOT NULL
    );`)
	if err != nil {
		return err
	}
	err = db.Exec(`CREATE INDEX idx_cap ON entry(cap);`)
	return err
}

func Record(db *genji.DB, alpha obj.AlphaMessage) error {
	return db.Exec(`INSERT INTO entry ( timestamp, cap, message ) VALUES ( ?, ?, ? )`, alpha.Timestamp.Unix(), alpha.CapCode, alpha.Message)
}
