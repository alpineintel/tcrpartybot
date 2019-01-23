package models

import (
	"database/sql"
	"math/big"
	"time"
)

// FaucetHit is a record of each time a user hits the faucet
type FaucetHit struct {
	ID        int64      `db:"id"`
	AccountID int64      `db:"account_id"`
	Amount    string     `db:"amount"`
	Timestamp *time.Time `db:"created_at"`
}

// RecordFaucetHit will create a new hit in the database for the given user
func RecordFaucetHit(accountID int64, amount *big.Int) error {
	db := GetDBSession()

	_, err := db.Exec(`
		INSERT INTO faucet_hits (
			account_id,
			amount
		) VALUES($1, $2)
	`, accountID, amount.String())

	return err
}

// LatestFaucetHit returns the record for the latest hit from the given user
func LatestFaucetHit(accountID int64) (*FaucetHit, error) {
	db := GetDBSession()

	hit := &FaucetHit{}
	err := db.Get(hit, `
		SELECT * FROM
			faucet_hits
		WHERE
			account_id=$1
		ORDER BY created_at
		LIMIT 1
	`, accountID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return hit, nil
}
