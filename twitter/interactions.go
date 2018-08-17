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
