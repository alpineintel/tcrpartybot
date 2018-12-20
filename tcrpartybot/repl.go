package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/events"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	HelpString = `Welcome to the TCR Party REPL! Available commands:
	dm [from handle, w/o @] [message]      - Simulates a Twitter DM
	mention [from handle, w/o @] [message] - Simulates a Twitter mention
	auth-vip                               - Begins auth bot auth flow
	auth-party                             - Begins retweet bot auth flow
	send-dm [to handle, w/o @] [message]   - Sends DM to a user from VIP bot
	create-webhook                         - Creates a webhook for listening on DMs
	distribute                             - Distributes tokens to all pre-registered accounts
	deploy-wallet                          - Calls the MultisigWalletFactory contract [for debugging]`
)

func beginRepl(eventChan chan<- *events.TwitterEvent, errChan chan<- error) {
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

		case "deploy-wallet":
			deployWallet(errChan)
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

			eventChan <- &events.TwitterEvent{
				SourceHandle: args[0],
				Message:      strings.Join(args[1:], " "),
				EventType:    events.TwitterEventTypeDM,
				Time:         time.Now(),
			}
			break

		case "mention":
			if argc < 2 {
				errChan <- errors.New("Invalid number of arguments for command mention")
				continue
			}

			eventChan <- &events.TwitterEvent{
				SourceHandle: args[0],
				Message:      strings.Join(args[1:], " "),
				EventType:    events.TwitterEventTypeMention,
				Time:         time.Now(),
			}
			break
		case "distribute":
			distributeTokens(errChan)
		}
	}
}