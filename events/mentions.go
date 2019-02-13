package events

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	welcomeMsg         = "Welcome to the party! Before we put you on the VIP list, we need to make sure you're a human. Here's an easy question for you: %s"
	notFollowingTweet  = "@%s In order for us to get started I'll need you to follow me, otherwise we can't interact via DM."
	notRegisteredTweet = "@%s You'll need to register in order to join the party. Follow me and tweet out \"hey @TCRPartyVIP let's party\" to get started"
)

func processFollow(event *TwitterEvent, errChan chan<- error) {
	// If they are following us but already have an un-verified account it means
	// that they've already sent us a "let's party" tweet. Since we can now DM
	// them we can also kick off the verification process
	account, err := models.FindAccountByTwitterID(event.SourceID)
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	if account != nil {
		// Did the unfollow us before completing the last challenge?
		activeChallenge, err := models.FindIncompleteChallenge(account.ID)
		if err != nil {
			errChan <- err
			return
		}

		// Nope, let's find them their next unsent challenge
		if activeChallenge == nil {
			activeChallenge, err = models.FindUnsentChallenge(account.ID)
			if err != nil {
				errChan <- err
				return
			}
		}

		// Still no challenge question? Okay let's abort...
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
	if strings.ToLower(event.SourceHandle) == strings.ToLower(os.Getenv("PARTY_BOT_HANDLE")) {
		return
	} else if strings.ToLower(event.SourceHandle) == strings.ToLower(os.Getenv("VIP_BOT_HANDLE")) {
		return
	}

	log.Printf("\nReceived mention from %s [%d]: %s", event.SourceHandle, event.SourceID, event.Message)
	err := models.CreateMentionAnalyticsEvent(event.SourceID, event.ObjectID, event.Message)
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	lower := strings.ToLower(event.Message)

	account, err := models.FindAccountByHandle(event.SourceHandle)
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	} else if account == nil {
		// No account. Do they want to regiser?
		if strings.Contains(lower, " party") {
			processRegistration(event, errChan)
			return
		}

		// No idea what they're doing, let's tweet at them
		tweet := fmt.Sprintf(notRegisteredTweet, event.SourceHandle)
		if err = twitter.SendTweet(twitter.VIPBotHandle, tweet); err != nil {
			errChan <- errors.Wrap(err)
		}
		return
	}

	// Define a helper function that will be passed around below
	sendDM := generateSendDM(account, errChan)

	// Are they trying to challenge or nominate?
	nominateMatcher := regexp.MustCompile("nominate @" + twitterHandleRegex)
	challengeMatcher := regexp.MustCompile("challenge @" + twitterHandleRegex)
	voteMatcher := regexp.MustCompile(`vote @` + twitterHandleRegex + ` (keep|kick) ?([0-9]*)`)

	if nominateMatcher.MatchString(lower) {
		matches := nominateMatcher.FindStringSubmatch(lower)
		if len(matches) < 2 {
			errChan <- errors.Errorf("could not parse mention nominee %s", matches)
			return
		}

		args := []string{"nominate", matches[1]}
		err = handleNomination(account, args, sendDM)
	} else if challengeMatcher.MatchString(lower) {
		matches := challengeMatcher.FindStringSubmatch(lower)
		if len(matches) < 2 {
			errChan <- errors.Errorf("could not parse mention challenge %s", matches)
			return
		}

		args := []string{"challenge", matches[1]}
		err = handleChallenge(account, args, sendDM)
	} else if voteMatcher.MatchString(lower) {
		matches := voteMatcher.FindStringSubmatch(lower)
		if len(matches) < 3 {
			errChan <- errors.Errorf("could not parse vote challenge %s", matches)
			return
		}

		args := []string{"vote", matches[1], matches[2]}
		if len(matches) == 4 {
			args = append(args, matches[3])
		}

		err = handleVote(account, args, sendDM)
	}

	if err != nil {
		errChan <- errors.Wrap(err)
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
