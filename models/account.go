package models

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

// Account represents a twitter account that interacts with our bot. They may
// or may not be a member on the TCR or a token holder.
type Account struct {
	ID                            int64           `db:"id"`
	TwitterID                     int64           `db:"twitter_id"`
	TwitterHandle                 string          `db:"twitter_handle"`
	MultisigAddress               *sql.NullString `db:"multisig_address"`
	MultisigFactoryIdentifier     *sql.NullInt64  `db:"multisig_factory_identifier"`
	PassedRegistrationChallengeAt *time.Time      `db:"passed_registration_challenge_at"`
	CreatedAt                     *time.Time      `db:"created_at"`
	LastDMAt                      *time.Time      `db:"last_dm_at"`

	// ActivatedAt is used for the pre-registration phase of the party. After
	// pre-registration it can be ignored
	ActivatedAt *time.Time `db:"activated_at"`
}

// CreateAccount creates a new account record in the database from the given
// account struct.
func CreateAccount(account *Account) error {
	db := GetDBSession()

	var id int64
	err := db.QueryRow(`
		INSERT INTO accounts (
			twitter_handle,
			twitter_id
		) VALUES($1, $2)
		RETURNING id
	`, account.TwitterHandle, account.TwitterID).Scan(&id)

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

// FindAccountByID searches for an account based on a user ID (not twitter ID)
func FindAccountByID(id int64) (*Account, error) {
	db := GetDBSession()

	account := Account{}
	err := db.Get(&account, "SELECT * FROM accounts WHERE id=$1", id)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &account, nil
}

// FindAccountByTwitterID searches for a given account based on the provided ID or
// returns nil if it cannot be found
func FindAccountByTwitterID(id int64) (*Account, error) {
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

// FindAccountByMultisigAddress returns an account based on the provided multisig address
func FindAccountByMultisigAddress(address string) (*Account, error) {
	db := GetDBSession()

	account := Account{}
	err := db.Get(&account, "SELECT * FROM accounts WHERE multisig_address=$1", address)

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

// HasCompletedChallenges returns true if this account has completed all of the
// required verification challenges (and therefore has a verified account)
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
			passed_registration_challenge_at = :passed_registration_challenge_at,
			activated_at = :activated_at,
			last_dm_at = :last_dm_at
		WHERE id = :id
	`, a)

	return err
}

// MarkRegistered updates the passed_registration_challenge_at column with the
// current timestamp
func (a *Account) MarkRegistered() error {
	now := time.Now().UTC()
	a.PassedRegistrationChallengeAt = &now

	return a.Save()
}

// MarkActivated updates the activated_at column with the current timestamp
func (a *Account) MarkActivated() error {
	now := time.Now().UTC()
	a.ActivatedAt = &now

	return a.Save()
}

// UpdateLastDMAt sets last_dm_at to the current timestamp and updates the db
func (a *Account) UpdateLastDMAt() error {
	now := time.Now().UTC()
	a.LastDMAt = &now

	return a.Save()
}

// SetMultisigFactoryIdentifier updates the user's account with the given
// multisig factory identifier
func (a *Account) SetMultisigFactoryIdentifier(identifier int64) error {
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

// Destroy will delete all records of the account from the db
func (a *Account) Destroy() error {
	db := GetDBSession()

	tx := db.MustBegin()

	_, err := tx.NamedExec("DELETE FROM registration_challenges WHERE account_id = :id", a)
	if err != nil {
		return err
	}

	_, err = tx.NamedExec("DELETE FROM accounts WHERE id = :id", a)
	if err != nil {
		return err
	}

	return tx.Commit()
}
