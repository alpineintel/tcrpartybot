package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"log"
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

// CreateWebhook creates a new webhook for the given token, allowing us to
// receive notifications for new DMs. This should only be used on the
// TCRPartyVIP bot.
func CreateWebhook() (string, error) {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return "", err
	}

	params := &twitter.AccountActivityRegisterWebhookParams{
		EnvName: "dev",
		URL:     "https://example.com/test",
	}
	webhook, _, err := client.AccountActivity.RegisterWebhook(params)

	if err != nil {
		return "", err
	}

	return webhook.ID, nil
}
