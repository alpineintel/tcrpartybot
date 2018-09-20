package events

import (
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

func processFollow(event *Event, errChan chan<- error) {
	err := twitter.Follow(event.SourceID)
	if err != nil {
		errChan <- err
	}
}
