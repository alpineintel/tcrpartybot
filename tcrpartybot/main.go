package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/tokenfoundry/tcrpartybot/events"
	"log"
	"os"
	"time"
)

/*
* polls: [from blockchain]
*	active challenges
* accounts:
*	twitter_handle
*	eth_address
*	private_key [not on blockchain]
*
 */

type TCRContract interface {
	Apply(applicant string, nominee string, amount uint) (string, error)
	Vote(voter string, nominee string, accept bool) error
	Challenge(challengee string, challenger string) error
	GetExpiry(nominee string) (time.Time, error)
}

type TwitterCredentials struct {
	ConsumerKey    string
	ConsumerSecret string
}

func listenToTwitter(twitterCredentials *TwitterCredentials, eventChan chan<- *events.Event, errorChan chan<- error) {
	for {
		// Initiate the DM channel
		eventChan <- &events.Event{
			EventType:    events.EventTypeMention,
			Time:         time.Now(),
			SourceHandle: "stevenleeg",
			Message:      "Let's party, @TCRPartyVIP",
		}
		time.Sleep(1 * time.Second)

		// Dummy response to party challenge
		eventChan <- &events.Event{
			EventType:    events.EventTypeDM,
			Time:         time.Now(),
			SourceHandle: "stevenleeg",
			Message:      "5",
		}
		time.Sleep(1 * time.Second)

		// Send another dummy response
		eventChan <- &events.Event{
			EventType:    events.EventTypeDM,
			Time:         time.Now(),
			SourceHandle: "stevenleeg",
			Message:      "7",
		}
		time.Sleep(1 * time.Second)
	}
}

func logErrors(errorChan <-chan error) {
	for {
		log.Println(<-errorChan)
	}
}

func main() {
	twitterCredentials := &TwitterCredentials{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	eventChan := make(chan *events.Event)
	errorChan := make(chan error)

	go listenToTwitter(twitterCredentials, eventChan, errorChan)
	go events.ProcessEvents(eventChan, errorChan)

	for range eventChan {
	}
}
