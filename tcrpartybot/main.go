package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gitlab.com/alpinefresh/tcrpartybot/api"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/events"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

func main() {
	// Some pre-boot config
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	godotenv.Load()

	// Start the db connection pool
	models.GetDBSession()

	// Check to see if we have credentials for the two twitter handles
	_, err := models.FindOAuthTokenByHandle(os.Getenv("PARTY_BOT_HANDLE"))
	if err != nil {
		log.Printf("Credentials for party bot not found. Please authenticate!")
	}

	_, err = models.FindOAuthTokenByHandle(os.Getenv("VIP_BOT_HANDLE"))
	if err != nil {
		log.Printf("Credentials for vip bot not found. Please authenticate!")
	}

	errChan := make(chan error)

	// Listen for and process any incoming twitter events
	go twitter.MonitorRatelimit()

	twitterEventChan := make(chan *events.TwitterEvent)
	go api.StartServer(twitterEventChan, errChan)
	go events.ProcessTwitterEvents(twitterEventChan, errChan)

	// Look for any existing applications/challenges that may need to be updated
	go events.ScheduleUpdates(errChan)

	// Start listening for relevant events on the blockchain
	ethEvents := make(chan *events.ETHEvent)
	twitterETHChan := make(chan *events.ETHEvent, 10)
	processETHChan := make(chan *events.ETHEvent, 10)

	go func() {
		// Fan events out to each listner
		for event := range ethEvents {
			twitterETHChan <- event
			processETHChan <- event
		}
	}()

	go events.StartBotListener(ethEvents, errChan)
	go events.ListenAndRetweet(twitterETHChan, errChan)
	go events.ProcessETHEvents(processETHChan, errChan)

	startRepl := flag.Bool("repl", false, "Starts the debug REPL")
	flag.Parse()
	if *startRepl {
		go errors.LogErrors(errChan)
		beginRepl(twitterEventChan, errChan)
	} else {
		errors.LogErrors(errChan)
	}
}
