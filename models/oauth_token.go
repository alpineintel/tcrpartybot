package models

import (
	"strings"
	"time"
)

type OAuthToken struct {
	ID               int64     `db:"id"`
	TwitterHandle    string    `db:"twitter_handle"`
	OAuthToken       string    `db:"oauth_token"`
	OAuthTokenSecret string    `db:"oauth_token_secret"`
	CreatedAt        time.Time `db:"created_at"`
}

func CreateOAuthToken(token *OAuthToken) error {
	db := GetDBSession()

	result := db.MustExec(`
		INSERT INTO oauth_tokens (
			twitter_handle,
			oauth_token,
			oauth_token_secret
		) VALUES($1, $2, $3)
	`, token.TwitterHandle, token.OAuthToken, token.OAuthTokenSecret)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	token.ID = id
	return nil
}

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
