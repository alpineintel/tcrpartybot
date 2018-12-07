package twitter

import (
	"github.com/stevenleeg/go-twitter/twitter"
	"log"
	"os"
	"strconv"
)

// FilterTweets will begin filtering tweets and outputting them to the returned
// channel
func FilterTweets(sinceID int64, twitterHandles []string) (*twitter.Stream, <-chan *twitter.Tweet, error) {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Filtering tweets by:")
	twitterIDs := make([]string, len(twitterHandles))
	for _, handle := range twitterHandles {
		user, _, err := client.Users.Show(&twitter.UserShowParams{
			ScreenName: handle,
		})

		if err != nil {
			return nil, nil, err
		}

		twitterIDs = append(twitterIDs, strconv.FormatInt(user.ID, 10))
		log.Printf("\t%s (%d)", user.ScreenName, user.ID)
	}

	params := &twitter.StreamFilterParams{
		Follow: twitterIDs,
	}

	stream, err := client.Streams.Filter(params)
	if err != nil {
		return nil, nil, err
	}

	// Start listening on a separate goroutine
	tweetChan := make(chan *twitter.Tweet)
	go func() {
		demux := twitter.NewSwitchDemux()
		demux.Tweet = func(tweet *twitter.Tweet) {
			tweetChan <- tweet
		}

		defer close(tweetChan)
		demux.HandleChan(stream.Messages)
	}()

	return stream, tweetChan, nil
}

// Retweet creates a new RT of the given tweet ID
func Retweet(handle string, tweetID int64) error {
	client, _, err := GetClientFromHandle(handle)
	if err != nil {
		return nil
	}

	log.Printf("Retweeting from %s: %d", handle, tweetID)
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return nil
	}
	_, _, err = client.Statuses.Retweet(tweetID, nil)
	return err
}

// SendTweet sends a new tweet from the given handle constant
func SendTweet(handle string, message string) error {
	client, _, err := GetClientFromHandle(handle)
	if err != nil {
		return nil
	}

	log.Printf("Tweeting from %s: %s", handle, message)
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return nil
	}

	_, _, err = client.Statuses.Update(message, nil)
	return err
}

// SendDM sends a direct message from the VIP party bot to the specified handle
func SendDM(recipientID int64, message string) error {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return err
	}

	log.Printf("Sending DM to %d: %s", recipientID, message)
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return nil
	}

	_, _, err = client.DirectMessages.EventsCreate(&twitter.DirectMessageEventsCreateParams{
		RecipientID: strconv.FormatInt(recipientID, 10),
		Text:        message,
	})

	if err != nil {
		log.Printf("Failed sending DM to %d: %s", recipientID, message)
		return err
	}

	return nil
}

// Follow will create a new friendship with the given user ID
func Follow(userID int64) error {
	client, _, err := GetClientFromHandle(VIPBotHandle)

	follow := true
	params := &twitter.FriendshipCreateParams{
		UserID: userID,
		Follow: &follow,
	}
	_, _, err = client.Friendships.Create(params)
	return err
}

// CreateWebhook creates a new webhook and subscribes it to the user, allowing
// us to receive notifications for new DMs. This should only be used on the
// TCRPartyVIP bot.
func CreateWebhook() (string, error) {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return "", err
	}

	webhookParams := &twitter.AccountActivityRegisterWebhookParams{
		EnvName: os.Getenv("TWITTER_ENV"),
		URL:     os.Getenv("BASE_URL") + "/webhooks/twitter",
	}
	webhook, _, err := client.AccountActivity.RegisterWebhook(webhookParams)

	if err != nil {
		return "", err
	}

	return webhook.ID, nil
}

// CreateSubscription subscribes the current webhook to the given user
func CreateSubscription() error {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return err
	}

	subParams := &twitter.AccountActivityCreateSubscriptionParams{
		EnvName: os.Getenv("TWITTER_ENV"),
	}
	_, err = client.AccountActivity.CreateSubscription(subParams)

	return err
}
