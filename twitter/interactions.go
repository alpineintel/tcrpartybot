package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"os"
	"strconv"
)

// SendDM sends a direct message from the VIP party bot to the specified handle
func SendDM(recipientID int64, message string) error {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return err
	}

	_, _, err = client.DirectMessages.EventsCreate(&twitter.DirectMessageEventsCreateParams{
		RecipientID: strconv.FormatInt(recipientID, 10),
		Text:        message,
	})

	if err != nil {
		log.Printf("Failed sending DM to %d: %s", recipientID, message)
		return err
	}

	log.Printf("Sent DM to %d: %s", recipientID, message)
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
