package models

import (
	"database/sql"
	"time"
)

// Account represents a twitter account that interacts with our bot. They may
// or may not be a member on the TCR or a token holder.
type Account struct {
	ID                            int64      `db:"id"`
	TwitterID                     int64      `db:"twitter_id"`
	TwitterHandle                 string     `db:"twitter_handle"`
	ETHAddress                    string     `db:"eth_address"`
	ETHPrivateKey                 string     `db:"eth_private_key"`
	PassedRegistrationChallengeAt *time.Time `db:"passed_registration_challenge_at"`
	CreatedAt                     *time.Time `db:"created_at"`
}

func CreateAccount(account *Account) error {
	db := GetDBSession()

	var id int64
	err := db.QueryRow(`
		INSERT INTO accounts (
			twitter_handle,
			twitter_id,
			eth_address,
			eth_private_key
		) VALUES($1, $2, $3, $4)
		RETURNING id
	`, account.TwitterHandle, account.TwitterID, account.ETHAddress, account.ETHPrivateKey).Scan(&id)

	if err != nil {
		return err
	}

	account.ID = id
	return nil
}

// FindAccountByHandle searches for a given account based on its handle or
// returns nil if it cannot be found
func FindAccountByHandle(handle string) (*Account, error) {
	db := GetDBSession()

	account := Account{}
	err := db.Get(&account, "SELECT * FROM accounts WHERE twitter_handle=$1", handle)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &account, nil
}

// FindAccountByID searches for a given account based on the provided ID or
// returns nil if it cannot be found
func FindAccountByID(id int64) (*Account, error) {
	db := GetDBSession()

	account := Account{}
	err := db.Get(&account, "SELECT * FROM accounts WHERE twitter_id=$1", id)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &account, nil
}

func AccountHasCompletedChallenges(accountId int64) (bool, error) {
	db := GetDBSession()

	// A user has completed the challenge phase if there are X registration
	// challenges with non-nil completed_at columns
	var count int
	err := db.Get(&count, `
		SELECT COUNT(*)
		FROM registration_challenges
		WHERE
			account_id = $1 AND
			completed_at IS NOT NULL
	`, accountId)

	if err != nil {
		return false, err
	}

	return count == RegistrationChallengeCount, nil
}

func MarkAccountRegistered(accountId int64) error {
	db := GetDBSession()

	now := time.Now()
	_, err := db.Exec(`
		UPDATE accounts
		SET passed_registration_challenge_at = $1
		WHERE id = $2
	`, &now, accountId)

	return err
}
