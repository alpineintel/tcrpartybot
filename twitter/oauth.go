package twitter

import (
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
)

const (
	OUT_OF_BAND = "oob"
)

type TwitterCredentials struct {
	ConsumerKey    string
	ConsumerSecret string
}

func getOAuthConfiguration(credentials *TwitterCredentials) *oauth1.Config {
	return &oauth1.Config{
		ConsumerKey:    credentials.ConsumerKey,
		ConsumerSecret: credentials.ConsumerSecret,
		Endpoint:       twitter.AuthorizeEndpoint,
		CallbackURL:    OUT_OF_BAND,
	}

}

func GetOAuthURL(credentials *TwitterCredentials) (string, error) {
	conf := getOAuthConfiguration(credentials)

	requestToken, _, err := conf.RequestToken()
	if err != nil {
		return "", err
	}

	authorizationURL, err := conf.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}

	return authorizationURL.String(), nil
}
