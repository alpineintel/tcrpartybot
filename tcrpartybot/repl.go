package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/events"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	helpString = `Welcome to the TCR Party REPL! Available commands:
	dm [from handle, w/o @] [message]      - Simulates a Twitter DM
	mention [from handle, w/o @] [message] - Simulates a Twitter mention
	mention-id [from id]  [message]        - Simulates a Twitter mention using an ID
	follow [id]                            - Simulates a follow from the given Twitter ID
	send-dm [to handle, w/o @] [message]   - Sends DM to a user from VIP bot
	rm [handle]                            - Deletes all records of the account associated with this Twitter handle
	distribute                             - Distributes tokens to all pre-registered accounts
	withdraw [listing handle]              - Calls the withdraw method for the giving listing
	sb [block number]                      - Sets the last block number to the provided number
	deploy-wallet                          - Calls the MultisigWalletFactory contract [for debugging]`
)

func beginRepl(eventChan chan<- *events.TwitterEvent, errChan chan<- error) {
	fmt.Print(helpString)

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
		case "deploy-wallet":
			deployWallet(errChan)
			break

		case "follow":
			if argc < 1 {
				errChan <- errors.New("Invalid number of arguments for command follow")
				continue
			}

			twitterID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				errChan <- err
				continue
			}

			eventChan <- &events.TwitterEvent{
				SourceID:  twitterID,
				EventType: events.TwitterEventTypeFollow,
				Time:      time.Now().UTC(),
			}
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
				Time:         time.Now().UTC(),
			}
			break

		case "sb":
			if argc < 1 {
				errChan <- errors.New("Invalid number of arguments for command sb")
				continue
			}

			models.SetKey(models.LatestSyncedBlockKey, args[0])
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
				Time:         time.Now().UTC(),
			}
			break

		case "rm":
			if argc < 1 {
				errChan <- errors.New("Invalid number of arguments for command rm")
				continue
			}

			err := deleteAccount(args[0])
			if err != nil {
				errChan <- err
			} else {
				log.Println("Account deleted")
			}
			break

		case "mention-id":
			if argc < 3 {
				errChan <- errors.New("Invalid number of arguments for command dm-id")
				continue
			}

			twitterID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				errChan <- err
				continue
			}

			eventChan <- &events.TwitterEvent{
				SourceHandle: args[1],
				SourceID:     twitterID,
				Message:      strings.Join(args[2:], " "),
				EventType:    events.TwitterEventTypeMention,
				Time:         time.Now().UTC(),
			}
			break

		case "withdraw":
			if argc < 1 {
				errChan <- errors.New("Invalid number of arguments for command withdraw")
				continue
			}

			listing, err := contracts.GetListingFromHandle(args[0])
			if err != nil {
				errChan <- err
				continue
			} else if listing == nil {
				log.Printf("No listing for %s", args[0])
				continue
			}

			// Call the withdraw method from their multisig contract
			subtract := contracts.GetAtomicTokenAmount(500)
			amt := listing.UnstakedDeposit.Sub(listing.UnstakedDeposit, subtract)
			contracts.Withdraw(args[0], amt)
			break

		case "distribute":
			distributeTokens(errChan)
		}
	}
}
