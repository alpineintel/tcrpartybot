package twitter

import (
	"github.com/dghubble/oauth1"
	twitterOAuth "github.com/dghubble/oauth1/twitter"
	"github.com/stevenleeg/go-twitter/twitter"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"log"
	"os"
)

const (
	// OutOfBand is used for setting the callback URL on the twitter OAuth
	// endpoint
	OutOfBand = "oob"

	// VIPBotHandle is a placeholder which is resolved to the real handle via
	// an environment variable.
	VIPBotHandle = "vip"
	// PartyBotHandle is a placeholder which is resolved to the real handle via
	// an environment variable.
	PartyBotHandle = "party"
)

// OAuthRequest collects all required information for completing the
// OAuth1 authentication flow
type OAuthRequest struct {
	Handle       string
	PIN          string
	RequestToken string
}

// GetClientFromHandle returns a twitter client given a handle constant (ie vip
// or party) or an explicit twitter handle.
func GetClientFromHandle(handle string) (*twitter.Client, *models.OAuthToken, error) {
	if handle == VIPBotHandle {
		handle = os.Getenv("VIP_BOT_HANDLE")
	} else if handle == PartyBotHandle {
		handle = os.Getenv("PARTY_BOT_HANDLE")
	}

	oauthToken, err := models.FindOAuthTokenByHandle(handle)
	if err != nil {
		log.Printf("Could not find OAuth token for %s", handle)
		return nil, nil, err
	}

	conf := getOAuthConfiguration()
	token := oauth1.NewToken(oauthToken.OAuthToken, oauthToken.OAuthTokenSecret)
	httpClient := conf.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient), oauthToken, nil
}

// GetClientFromToken returns a twitter client given an OAuth token
func GetClientFromToken(oauthToken *models.OAuthToken) *twitter.Client {
	conf := getOAuthConfiguration()
	token := oauth1.NewToken(oauthToken.OAuthToken, oauthToken.OAuthTokenSecret)
	httpClient := conf.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

// GetOAuthURL returns a URL that can be used to authenticate a user via the
// OAuth 1 API
func (request *OAuthRequest) GetOAuthURL() (string, error) {
	conf := getOAuthConfiguration()

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

// ReceivePIN completes the OAuth dance by retrieving user info. Note that it
// requires the PIN to be set on the OAuthRequest struct.
func (request *OAuthRequest) ReceivePIN() error {
	conf := getOAuthConfiguration()
	accessToken, accessSecret, err := conf.AccessToken(
		request.RequestToken,
		"",
		request.PIN,
	)

	if err != nil {
		return err
	}

	// Save the OAuth credentials in the database
	token := &models.OAuthToken{
		TwitterHandle:    request.Handle,
		OAuthToken:       accessToken,
		OAuthTokenSecret: accessSecret,
	}

	// Fetch the user ID
	client := GetClientFromToken(token)
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(false),
	}
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		log.Println("Could fetch user data")
		return err
	}

	token.TwitterID = user.ID

	err = models.CreateOAuthToken(token)
	if err != nil {
		return err
	}

	return nil
}

func getOAuthConfiguration() *oauth1.Config {
	return &oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
		Endpoint:       twitterOAuth.AuthorizeEndpoint,
		CallbackURL:    OutOfBand,
	}
}
