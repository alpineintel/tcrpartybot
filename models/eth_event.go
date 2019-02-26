package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type ETHEvent struct {
	ID          int64          `db:"id"`
	EventType   string         `db:"event_type"`
	Data        types.JSONText `db:"data"`
	BlockNumber uint64         `db:"block_number"`
	TxHash      string         `db:"tx_hash"`
	TxIndex     uint           `db:"tx_index"`
	LogIndex    uint           `db:"log_index"`
	CreatedAt   *time.Time     `db:"created_at"`
}

// CreateETHEvent creates a new event given an Ethereum event name and its
// associated data
func CreateETHEvent(event *ETHEvent, data interface{}) error {
	db := GetDBSession()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	event.Data = types.JSONText(string(bytes))

	var id int64
	err = db.QueryRow(`
		INSERT INTO eth_events (
			event_type,
			block_number,
			data,
			tx_hash,
			tx_index,
			log_index,
			created_at
		) VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, event.EventType, event.BlockNumber, event.Data, event.TxHash, event.TxIndex, event.LogIndex, event.CreatedAt).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

// FindETHEventByID returns an ETH event given its ID
func FindETHEventByID(id int64) (*ETHEvent, error) {
	db := GetDBSession()

	ethEvent := &ETHEvent{}
	err := db.Get(ethEvent, "SELECT * FROM eth_events WHERE id=$1", id)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return ethEvent, nil
}

// FindETHEventsSinceID returns all events that have occured with IDs greater
// than the provided number. Since initial syncs have the potential of loading
// thousands of events, this function will return at most 500 (this number was
// selected arbitrarily and can be replaced based on performance) and return
// true on the moreAvailable value, signaling that the user should call the
// function again to retrieve the next batch of results.
func FindETHEventsSinceID(id int64) (events []*ETHEvent, moreAvailable bool, err error) {
	db := GetDBSession()
	err = db.Select(&events, "SELECT * FROM eth_events WHERE id > $1 ORDER BY id ASC LIMIT 500", id)

	if err != nil {
		return nil, false, err
	}

	return events, len(events) == 500, nil
}
