package models

import (
	"strings"
	"time"
)

// OAuthToken stores the data required to control the two bots which act as the
// frontend for our TCR. This table should only have two rows, one for the TCR
// Bot and one for the VIP bot.
type OAuthToken struct {
	ID               int64     `db:"id"`
	TwitterHandle    string    `db:"twitter_handle"`
	TwitterID        int64     `db:"twitter_id"`
	OAuthToken       string    `db:"oauth_token"`
	OAuthTokenSecret string    `db:"oauth_token_secret"`
	CreatedAt        time.Time `db:"created_at"`
}

// CreateOAuthToken inserts a new OAuth token into the database
func CreateOAuthToken(token *OAuthToken) error {
	db := GetDBSession()

	result := db.MustExec(`
		INSERT INTO oauth_tokens (
			twitter_handle,
			oauth_token,
			oauth_token_secret,
			twitter_id
		) VALUES($1, $2, $3, $4)
	`, token.TwitterHandle, token.OAuthToken, token.OAuthTokenSecret, token.TwitterID)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	token.ID = id
	return nil
}

// FindOAuthTokenByHandle returns an OAuth token given a twitter handle
func FindOAuthTokenByHandle(handle string) (*OAuthToken, error) {
	db := GetDBSession()

	token := OAuthToken{}
	handle = strings.ToLower(handle)
	err := db.Get(&token, "SELECT * FROM oauth_tokens WHERE LOWER(twitter_handle)=$1", handle)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Save saves the OAuthToken to the database
func (token *OAuthToken) Save() error {
	db := GetDBSession()

	result := db.MustExec(`
		UPDATE oauth_tokens
		SET
			twitter_handle = $1,
			oauth_token = $2,
			oauth_token_secret = $3,
			twitter_id = $4
		WHERE
			id=$5
	`, token.TwitterHandle, token.OAuthToken, token.OAuthTokenSecret, token.TwitterID, token.ID)

	_, err := result.RowsAffected()
	return err
}
