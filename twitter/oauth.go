package twitter

import (
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/tokenfoundry/tcrpartybot/models"
)

const (
	OUT_OF_BAND = "oob"
)

type TwitterCredentials struct {
	ConsumerKey    string
	ConsumerSecret string
}

type TwitterOAuthRequest struct {
	Handle       string
	PIN          string
	RequestToken string
}

func getOAuthConfiguration(credentials *TwitterCredentials) *oauth1.Config {
	return &oauth1.Config{
		ConsumerKey:    credentials.ConsumerKey,
		ConsumerSecret: credentials.ConsumerSecret,
		Endpoint:       twitter.AuthorizeEndpoint,
		CallbackURL:    OUT_OF_BAND,
	}

}

func GetOAuthURL(credentials *TwitterCredentials, request *TwitterOAuthRequest) (string, error) {
	conf := getOAuthConfiguration(credentials)

	requestToken, _, err := conf.RequestToken()
	if err != nil {
		return "", err
	}

	authorizationURL, err := conf.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}

	request.RequestToken = requestToken

	return authorizationURL.String(), nil
}

func ReceivePIN(credentials *TwitterCredentials, request *TwitterOAuthRequest) error {
	conf := getOAuthConfiguration(credentials)
	accessToken, accessSecret, err := conf.AccessToken(
		request.RequestToken,
		"",
		request.PIN,
	)

	if err != nil {
		return err
	}

	token := &models.OAuthToken{
		TwitterHandle:    request.Handle,
		OAuthToken:       accessToken,
		OAuthTokenSecret: accessSecret,
	}

	return models.CreateOAuthToken(token)
}
