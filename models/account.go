package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

// Account represents a twitter account that interacts with our bot. They may
// or may not be a member on the TCR or a token holder.
type Account struct {
	ID                            int64           `db:"id"`
	TwitterID                     int64           `db:"twitter_id"`
	TwitterHandle                 string          `db:"twitter_handle"`
	ETHAddress                    string          `db:"eth_address"`
	ETHPrivateKey                 string          `db:"eth_private_key"`
	MultisigAddress               *sql.NullString `db:"multisig_address"`
	MultisigFactoryIdentifier     *sql.NullInt64  `db:"multisig_factory_identifier"`
	PassedRegistrationChallengeAt *time.Time      `db:"passed_registration_challenge_at"`
	CreatedAt                     *time.Time      `db:"created_at"`
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

// AllAccounts returns a list of all accounts currently in the database
func AllAccounts() (*sqlx.Rows, error) {
	db := GetDBSession()
	rows, err := db.Queryx("SELECT * FROM accounts")

	return rows, err
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

// FindAccountByMultisigFactoryIdentifier searches for a given account based on
// its factory identifier returns nil if it cannot be found
func FindAccountByMultisigFactoryIdentifier(identifier int64) (*Account, error) {
	db := GetDBSession()

	account := Account{}
	err := db.Get(&account, "SELECT * FROM accounts WHERE multisig_factory_identifier=$1", identifier)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &account, nil
}

func (a *Account) HasCompletedChallenges() (bool, error) {
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
	`, a.ID)

	if err != nil {
		return false, err
	}

	return count == RegistrationChallengeCount, nil
}

// Save updates the account record in the database and returns an error if it occurs
func (a *Account) Save() error {
	db := GetDBSession()

	_, err := db.NamedExec(`
		UPDATE accounts SET
			twitter_id = :twitter_id,
			twitter_handle = :twitter_handle,
			multisig_address = :multisig_address,
			multisig_factory_identifier = :multisig_factory_identifier,
			passed_registration_challenge_at = :passed_registration_challenge_at
		WHERE id = :id
	`, a)

	return err
}

// MarkRegistered updates the passed_registration_challenge_at column with the
// current timestamp
func (a *Account) MarkRegistered() error {
	now := time.Now()
	a.PassedRegistrationChallengeAt = &now

	return a.Save()
}

// SetMultisigFactoryIdentifier updates the user's account with the given
// multisig factory identifier
func (a *Account) SetMultisigFactoryIdentifier(identifier int64) error {
	fmt.Println(a.MultisigFactoryIdentifier)
	if a.MultisigFactoryIdentifier != nil && a.MultisigFactoryIdentifier.Valid {
		return errors.New("identifier can only be set once")
	}

	a.MultisigFactoryIdentifier = &sql.NullInt64{Valid: true, Int64: identifier}
	return a.Save()
}

// SetMultisigAddress updates the user's account with the given multisig wallet
// address
func (a *Account) SetMultisigAddress(address string) error {
	a.MultisigAddress = &sql.NullString{Valid: true, String: address}
	return a.Save()
}
