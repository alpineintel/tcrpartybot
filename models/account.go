package models

type Account struct {
	ID            int64
	TwitterHandle string `db:"twitter_handle"`
	ETHAddress    string `db:"eth_address"`
	ETHPrivateKey string `db:"eth_private_key"`
}

func CreateAccount(account *Account, errChan chan<- error) {
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
		errChan <- err
		return
	}

	account.ID = id
}
