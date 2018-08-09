package models

type Account struct {
	ID            int64  `db:"id"`
	TwitterHandle string `db:"twitter_handle"`
	ETHAddress    string `db:"eth_address"`
	ETHPrivateKey string `db:"eth_private_key"`
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

func FindAccountByHandle(handle string) *Account {
	db := GetDBSession()

	account := Account{}
	err := db.Get(&account, "SELECT * FROM accounts WHERE twitter_handle=$1", handle)

	if err != nil {
		return nil
	}

	return &account
}
