package events

import (
	"errors"
	"fmt"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
	"log"
	"os"
	"strings"
)

const (
	welcomeMsg        = "Welcome to the party! Before we put you on the VIP list, we need to make sure you're a human. Here's an easy question for you: %s"
	notFollowingTweet = "@%s Hey! In order for us to get started I'll need you to follow me, otherwise we can't interact via DM."
)

func processFollow(event *TwitterEvent, errChan chan<- error) {
	// Follow them back
	err := twitter.Follow(event.SourceID)
	if err != nil {
		errChan <- err
	}

	// If they are following us but already have an un-verified account it means
	// that they've already sent us a "let's party" tweet. Since we can now DM
	// them we can also kick off the verification process
	account, err := models.FindAccountByTwitterID(event.SourceID)
	if account != nil {
		// Did the unfollow us before completing the last challenge?
		activeChallenge, err := models.FindIncompleteChallenge(account.ID)
		if err != nil {
			errChan <- err
			return
		}

		if activeChallenge == nil {
			activeChallenge, err = models.FindUnsentChallenge(account.ID)
			if err != nil {
				errChan <- err
				return
			}
		}

		if activeChallenge == nil {
			log.Printf("Could not find active challenge for user %s, aborting follow DM response!", event.SourceHandle)
			return
		}

		err = twitter.SendDM(account.TwitterID, fmt.Sprintf(welcomeMsg, activeChallenge.Question))
		if err != nil {
			errChan <- err
			return
		}

		err = activeChallenge.MarkSent()
		if err != nil {
			errChan <- err
			return
		}
	}
}

func processMention(event *TwitterEvent, errChan chan<- error) {
	if strings.ToLower(event.SourceHandle) == os.Getenv("PARTY_BOT_HANDLE") {
		return
	} else if strings.ToLower(event.SourceHandle) == os.Getenv("VIP_BOT_HANDLE") {
		return
	}

	log.Printf("\nReceived mention from %s [%d]: %s", event.SourceHandle, event.SourceID, event.Message)

	// Filter based on let's party
	lower := strings.ToLower(event.Message)
	if strings.Contains(lower, "party") {
		processRegistration(event, errChan)
	}
}

func processRegistration(event *TwitterEvent, errChan chan<- error) {
	// If they already have an account we don't need to continue
	account, err := models.FindAccountByTwitterID(event.SourceID)
	if account != nil {
		return
	} else if err != nil {
		errChan <- err
		return
	}

	log.Printf("Creating account for %d", event.SourceID)

	account = &models.Account{
		TwitterHandle: event.SourceHandle,
		TwitterID:     event.SourceID,
	}
	err = models.CreateAccount(account)

	if err != nil {
		errChan <- err
	}

	// Generate three registration challenges for them
	questions, err := models.FetchRandomRegistrationQuestions(models.RegistrationChallengeCount)

	if err != nil {
		errChan <- err
		return
	}

	// Create a list of challenges for the new user to complete
	challenges := make([]*models.RegistrationChallenge, models.RegistrationChallengeCount)
	for i, question := range questions {
		challenges[i], err = models.CreateRegistrationChallenge(account, &question)

		if err != nil {
			errChan <- err
			return
		}
	}

	// If they're not currently following us we need to send them a mention
	// telling them to, otherwise we can't send a DM
	isFollowing, err := twitter.IsFollower(event.SourceID)
	if err != nil {
		errChan <- err
		return
	} else if !isFollowing {
		log.Printf("%s is not following, sending message.", account.TwitterHandle)
		text := fmt.Sprintf(notFollowingTweet, account.TwitterHandle)
		if err = twitter.SendTweet(twitter.VIPBotHandle, text); err != nil {
			errChan <- err
		}
		return
	}

	// Send them a direct message asking them for the answer to a challenge
	// question
	if len(questions) == 0 {
		errChan <- errors.New("Could not fetch registration question from db")
		return
	}

	firstChallenge := challenges[0]
	err = twitter.SendDM(account.TwitterID, fmt.Sprintf(welcomeMsg, questions[0].Question))
	if err != nil {
		errChan <- err
		return
	}

	err = firstChallenge.MarkSent()
	if err != nil {
		errChan <- err
	}
}
