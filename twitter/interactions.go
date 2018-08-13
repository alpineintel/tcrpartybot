package twitter

import (
	//"github.com/dghubble/go-twitter/twitter"
	"log"
)

func SendDM(handle string, message string) error {
	//client := GetClient(handle)
	log.Printf("Sent DM to %s: %s", handle, message)
	return nil
}
