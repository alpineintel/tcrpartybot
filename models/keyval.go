package models

import (
	"database/sql"
)

const (
	LatestSyncedBlockKey = "LatestSyncedBlock"
	LatestSyncedTweetKey = "LatestSyncedTweet"
)

type keyValueRow struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

// SetKey updates or inserts a key with the provided value
func SetKey(key string, value string) error {
	db := GetDBSession()

	result := db.MustExec(`
		INSERT INTO keyval_store (key, value) VALUES($1, $2)
		ON CONFLICT (key) DO UPDATE SET value=$2
	`, key, value)

	_, err := result.RowsAffected()
	return err
}

// GetKey fetches the given key's value from the database
func GetKey(key string) (string, error) {
	db := GetDBSession()

	row := &keyValueRow{}
	err := db.Get(row, "SELECT * FROM keyval_store WHERE key=$1 LIMIT 1", key)

	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	return row.Value, nil
}
