package models

import (
	"time"
)

type Account struct {
	ID                            int64      `db:"id"`
	TwitterHandle                 string     `db:"twitter_handle"`
	ETHAddress                    string     `db:"eth_address"`
	ETHPrivateKey                 string     `db:"eth_private_key"`
	PassedRegistrationChallengeAt *time.Time `db:"passed_registration_challenge_at"`
	CreatedAt                     *time.Time `db:"created_at"`
}

func CreateAccount(account *Account) error {
	db := GetDBSession()

	result := db.MustExec(`
		INSERT INTO accounts (
			twitter_handle,
			eth_address,
			eth_private_key
		) VALUES($1, $2, $3)
	`, account.TwitterHandle, account.ETHAddress, account.ETHPrivateKey)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	account.ID = id
	return nil
}

func FindAccountByHandle(handle string) (*Account, error) {
	db := GetDBSession()

	account := Account{}
	err := db.Get(&account, "SELECT * FROM accounts WHERE twitter_handle=$1", handle)

	if err != nil {
		return nil, err
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

	return count == REGISTRATION_CHALLENGE_COUNT, nil
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
