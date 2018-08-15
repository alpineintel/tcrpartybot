package models

import (
	"database/sql"
	"strings"
	"time"
)

type OAuthToken struct {
	ID               int64         `db:"id"`
	TwitterHandle    string        `db:"twitter_handle"`
	TwitterID        sql.NullInt64 `db:"twitter_id"`
	OAuthToken       string        `db:"oauth_token"`
	OAuthTokenSecret string        `db:"oauth_token_secret"`
	CreatedAt        time.Time     `db:"created_at"`
}

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
