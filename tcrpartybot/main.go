package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tokenfoundry/tcrpartybot/events"
	"github.com/tokenfoundry/tcrpartybot/models"
	"github.com/tokenfoundry/tcrpartybot/twitter"
	"log"
	"os"
	"strings"
	"time"
)

const (
	HELP_STRING = `Welcome to the TCR Party REPL! Available commands:
	dm [from handle, w/o @] [message]      - Simulates a Twitter DM
	mention [from handle, w/o @] [message] - Simulates a Twitter mention
	auth-vip                               - Begins auth bot auth flow
	auth-party                             - Begins retweet bot auth flow`
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

func logErrors(errorChan <-chan error) {
	for err := range errorChan {
		log.Printf("\n%s", err)
	}
}

func beginRepl(eventChan chan<- *events.Event, errChan chan<- error) {
	fmt.Print(HELP_STRING)

	twitterCredentials := &twitter.TwitterCredentials{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	for {
		// Give the other channels a chance to process and print a response
		time.Sleep(150 * time.Millisecond)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n(tcrparty)> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			errChan <- err
			return
		}

		trimmed := strings.TrimSpace(text)
		argv := strings.Split(trimmed, " ")
		command := argv[0]
		args := argv[1:]
		argc := len(args)

		switch command {
		case "auth-vip":
			request := &twitter.TwitterOAuthRequest{
				Handle: os.Getenv("VIP_BOT_HANDLE"),
			}

			url, err := twitter.GetOAuthURL(twitterCredentials, request)
			if err != nil {
				errChan <- err
				continue
			}

			fmt.Printf("Go to this URL to generate an access token:\n%s", url)
			fmt.Print("\nEnter PIN: ")

			_, err = fmt.Scanf("%s", &request.PIN)
			if err != nil {
				log.Println("Error receiving PIN")
				errChan <- err
				continue
			}

			err = twitter.ReceivePIN(twitterCredentials, request)
			if err != nil {
				log.Println("Error fetching token")
				errChan <- err
				continue
			}

			log.Println("Access token saved!")
			break
		case "dm":
			if argc < 2 {
				errChan <- errors.New("Invalid number of arguments for command dm")
				continue
			}

			eventChan <- &events.Event{
				SourceHandle: args[0],
				Message:      strings.Join(args[1:], " "),
				EventType:    events.EventTypeDM,
				Time:         time.Now(),
			}
			break

		case "mention":
			if argc < 2 {
				errChan <- errors.New("Invalid number of arguments for command mention")
				continue
			}

			eventChan <- &events.Event{
				SourceHandle: args[0],
				Message:      strings.Join(args[1:], " "),
				EventType:    events.EventTypeMention,
				Time:         time.Now(),
			}
			break
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not open .env file")
	}

	eventChan := make(chan *events.Event)
	errorChan := make(chan error)

	models.GetDBSession()

	go events.ProcessEvents(eventChan, errorChan)
	go logErrors(errorChan)

	beginRepl(eventChan, errorChan)
}
