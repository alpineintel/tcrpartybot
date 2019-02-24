package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"math/big"
)

// Balance represents an account's current balance at a given block
type Balance struct {
	ID            int64       `db:"id"`
	AccountID     int64       `db:"account_id"`
	ETHEventID    int64       `db:"eth_event_id"`
	PLCRBalance   TokenAmount `db:"plcr_balance"`
	WalletBalance TokenAmount `db:"wallet_balance"`
}

// TokenAmount represents a big.Int stored in a postgres VARCHAR column
type TokenAmount struct {
	Amt *big.Int
}

func (t TokenAmount) Value() (driver.Value, error) {
	return driver.Value(t.Amt.String()), nil
}

func (t *TokenAmount) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	default:
		return errors.New("incompatible type for TokenAmount")
	}

	t.Amt = new(big.Int)
	t.Amt.SetString(source, 10)
	return nil
}

// CreateBalance creates a new balance entry in the database
func CreateBalance(accountID, ethEventID int64, walletBalance, plcrBalance *big.Int) (*Balance, error) {
	db := GetDBSession()

	balance := &Balance{
		AccountID:     accountID,
		ETHEventID:    ethEventID,
		WalletBalance: TokenAmount{walletBalance},
		PLCRBalance:   TokenAmount{plcrBalance},
	}

	_, err := db.NamedExec(`
		INSERT INTO balances (
			account_id,
			eth_event_id,
			plcr_balance,
			wallet_balance
		) VALUES (
			:account_id,
			:eth_event_id,
			:plcr_balance,
			:wallet_balance
		)
	`, balance)

	return balance, err
}

// FindLatestUserBalance returns the most recent balance row for the given account ID
func FindLatestUserBalance(accountID int64) (*Balance, error) {
	db := GetDBSession()

	balance := &Balance{}
	err := db.Get(balance, "SELECT balances.* FROM balances JOIN eth_events ON balances.eth_event_id = eth_events.id WHERE account_id=$1 ORDER BY id DESC LIMIT 1", accountID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return &Balance{
			AccountID:     accountID,
			ETHEventID:    0,
			WalletBalance: TokenAmount{big.NewInt(0)},
			PLCRBalance:   TokenAmount{big.NewInt(0)},
		}, nil
	}

	return balance, nil
}

// FindLatestBalance finds the most recently created balance row. This is
// usually used for syncing with the blockchain, as this row will represent the
// latest block number that has been synced.
func FindLatestBalance() (*Balance, error) {
	db := GetDBSession()

	balance := Balance{}
	err := db.Get(&balance, "SELECT balances.* FROM balances JOIN eth_events ON balances.eth_event_id = eth_events.id ORDER BY id DESC LIMIT 1")

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &balance, nil
}

// AddToPLCRBalance adds the given amount to the user's latest PLCR balance,
// creates a new balance entry, and returns it.
func (b *Balance) AddToPLCRBalance(ethEventID int64, amount *big.Int) (*Balance, error) {
	currentBalance := new(big.Int)
	currentBalance.Set(b.PLCRBalance.Amt)
	newBalance := currentBalance.Add(currentBalance, amount)

	return CreateBalance(b.AccountID, ethEventID, b.WalletBalance.Amt, newBalance)
}

// AddToWalletBalance adds the given amount to the user's latest wallet balance,
// creates a new balance entry, and returns it.
func (b *Balance) AddToWalletBalance(ethEventID int64, amount *big.Int) (*Balance, error) {
	currentBalance := new(big.Int)
	currentBalance.Set(b.WalletBalance.Amt)
	newBalance := currentBalance.Add(currentBalance, amount)

	return CreateBalance(b.AccountID, ethEventID, newBalance, b.PLCRBalance.Amt)
}
