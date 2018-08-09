package twitter

import (
	"log"
)

func SendDirectMessage(handle string, message string) error {
	log.Printf("Sent direct message to %s: %s", handle, message)
	return nil
}
