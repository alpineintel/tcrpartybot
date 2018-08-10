package twitter

import (
	"log"
)

func SendDM(handle string, message string) error {
	log.Printf("Sent DM to %s: %s", handle, message)
	return nil
}
