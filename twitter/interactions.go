package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"os"
)

// SendDM sends a direct message from the VIP party bot to the specified handle
func SendDM(handle string, message string) error {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return err
	}

	_, _, err = client.DirectMessages.New(&twitter.DirectMessageNewParams{
		ScreenName: handle,
		Text:       message,
	})

	if err != nil {
		log.Printf("Failed sending DM to %s: %s", handle, message)
		return err
	}

	log.Printf("Sent DM to %s: %s", handle, message)
	return nil
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
		EnvName: "dev",
		URL:     os.Getenv("BASE_URL") + "/webhooks/twitter",
	}
	webhook, _, err := client.AccountActivity.RegisterWebhook(webhookParams)

	if err != nil {
		return "", err
	}

	return webhook.ID, nil
}

func CreateSubscription() error {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return err
	}

	subParams := &twitter.AccountActivityCreateSubscriptionParams{
		EnvName: "dev",
	}
	_, err = client.AccountActivity.CreateSubscription(subParams)

	return err
}
