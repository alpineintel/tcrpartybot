package models

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type ETHEvent struct {
	ID          int64          `db:"id"`
	EventType   string         `db:"event_type"`
	Data        types.JSONText `db:"additionals"`
	BlockNumber uint64         `db:"block_number"`
	CreatedAt   *time.Time     `db:"created_at"`
}

// CreateETHEvent creates a new event given an Ethereum event name and its
// associated data
func CreateETHEvent(eventType string, blockNumber uint64, timestamp *time.Time, data interface{}) error {
	db := GetDBSession()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	event := &ETHEvent{
		EventType: eventType,
		Data:      types.JSONText(string(bytes)),
		CreatedAt: timestamp,
	}

	var id int64
	err = db.QueryRow(`
		INSERT INTO eth_events (
			event_type,
			block_number,
			data,
			created_at
		) VALUES($1, $2, $3, $4)
		RETURNING id
	`, event.EventType, blockNumber, event.Data, event.CreatedAt).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}
