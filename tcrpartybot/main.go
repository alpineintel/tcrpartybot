package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"gitlab.com/alpinefresh/tcrpartybot/api"
	"gitlab.com/alpinefresh/tcrpartybot/events"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	HelpString = `Welcome to the TCR Party REPL! Available commands:
	dm [from handle, w/o @] [message]      - Simulates a Twitter DM
	mention [from handle, w/o @] [message] - Simulates a Twitter mention
	auth-vip                               - Begins auth bot auth flow
	auth-party                             - Begins retweet bot auth flow
	send-dm [to handle, w/o @] [message]   - Sends DM to a user from VIP bot
	create-webhook                         - Creates a webhook for listening on DMs`
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

func logErrors(errChan <-chan error) {
	for err := range errChan {
		log.Printf("\n%s", err)
	}
}

func authenticateHandle(handle string, errChan chan<- error) {
	request := &twitter.OAuthRequest{
		Handle: handle,
	}

	url, err := request.GetOAuthURL()
	if err != nil {
		errChan <- err
		return
	}

	fmt.Printf("Go to this URL to generate an access token:\n%s", url)
	fmt.Print("\nEnter PIN: ")

	_, err = fmt.Scanf("%s", &request.PIN)
	if err != nil {
		log.Println("Error receiving PIN")
		errChan <- err
		return
	}

	err = request.ReceivePIN()
	if err != nil {
		log.Println("Error fetching token")
		errChan <- err
		return
	}

	log.Println("Access token saved!")
}

func createWebhook(errChan chan<- error) {
	webhookID, err := models.GetKey("webhookID")
	if err != nil {
		errChan <- err
		return
	}

	// If we don't already have a webhook ID we should create it
	if webhookID == "" {
		id, err := twitter.CreateWebhook()
		if err != nil {
			errChan <- err
			return
		}

		log.Printf("Webhook %s created successfully", id)
		models.SetKey("webhookID", id)
	}

	// And subscribe to TCRPartyVIP's DMs
	if err := twitter.CreateSubscription(); err != nil {
		errChan <- err
		return
	}
	log.Printf("Subscription created successfully")
}

func beginRepl(eventChan chan<- *events.Event, errChan chan<- error) {
	fmt.Print(HelpString)

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
			authenticateHandle(os.Getenv("VIP_BOT_HANDLE"), errChan)
			break
		case "auth-party":
			authenticateHandle(os.Getenv("PARTY_BOT_HANDLE"), errChan)
			break

		case "create-webhook":
			createWebhook(errChan)
			break
		case "send-dm":
			if argc < 2 {
				errChan <- errors.New("Invalid number of arguments for command send-dm")
				continue
			}

			twitterID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				errChan <- err
				continue
			}

			err = twitter.SendDM(twitterID, strings.Join(args[1:], " "))
			if err != nil {
				errChan <- err
				continue
			}
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
	errChan := make(chan error)

	models.GetDBSession()

	_, err = models.FindOAuthTokenByHandle(os.Getenv("PARTY_BOT_HANDLE"))
	if err != nil {
		log.Printf("Credentials for party bot not found. Please authenticate!")
	}

	_, err = models.FindOAuthTokenByHandle(os.Getenv("VIP_BOT_HANDLE"))
	if err != nil {
		log.Printf("Credentials for vip bot not found. Please authenticate!")
	}

	go events.ProcessEvents(eventChan, errChan)
	go api.StartServer(eventChan, errChan)
	go logErrors(errChan)

	beginRepl(eventChan, errChan)
}
