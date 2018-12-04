package events

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

// ScheduleUpdates finds new/existing applications/challenge polls and
// schedules tasks to fire their on-chain updater functions once they have
// matured. Example: calling updateStatus on an application once the
// application period has passed.
func ScheduleUpdates(eventChan <-chan *ETHEvent, errChan chan<- error) {
	// First let's instantiate ourselves with existing challenges/applications
	// that need to be scheduled
	listings, err := contracts.GetAllListings()
	if err != nil {
		errChan <- err
		return
	}

	for _, listing := range listings {
		go scheduleListing(listing, errChan)
	}

	// Listen for any incoming applications and queue them up
	for {
		event := <-eventChan
		var err error

		switch event.EventType {
		case ETHEventNewTCRApplication:
			ethEvent, err := contracts.DecodeApplicationEvent(event.Topics, event.Data)
			if err != nil {
				errChan <- err
				continue
			}

			listing, err := contracts.GetListingFromHash(ethEvent.ListingHash)
			if err != nil {
				errChan <- err
				continue
			}

			go scheduleListing(listing, errChan)
			break
		}

		if err != nil {
			errChan <- err
		}
	}
}

func scheduleListing(application *contracts.RegistryListing, errChan chan<- error) {
	hasOpenChallenge := application.ChallengeID.Cmp(big.NewInt(0)) != 0

	if !application.Whitelisted && !hasOpenChallenge {
		// This listing hasn't been whitelisted yet and doesn't have an open
		// challenge. This means we'll need to schedule a updateStatus task
		expirationTime := time.Unix(application.ApplicationExpiry.Int64(), 0)
		if expirationTime.After(time.Now()) {
			time.Sleep(time.Until(expirationTime))
		}

		updateStatus(application, errChan)
	} else if hasOpenChallenge {
		// The listing has an open challenge, meaning we'll need to schedule
		// tasks to reveal any votes and update the status
		poll, err := contracts.GetPoll(application.ChallengeID)
		if err != nil {
			errChan <- err
			return
		}

		twitterHandle, err := contracts.GetListingDataFromHash(application.ListingHash)
		if err != nil {
			errChan <- err
			return
		}

		commitEndTime := time.Unix(poll.CommitEndDate.Int64(), 0)
		revealEndTime := time.Unix(poll.RevealEndDate.Int64(), 0)
		if commitEndTime.After(time.Now()) {
			// We haven't yet hit the commit time, so let's sleep until we do
			// and then reveal the vote
			log.Printf("[updater] Challenge @%s is in commit. Sleeping until %s", twitterHandle, commitEndTime.Format(time.UnixDate))
			time.Sleep(time.Until(commitEndTime) + (2 * time.Minute))
		}

		if revealEndTime.After(time.Now()) {
			reveal(application, errChan)
			log.Printf("[updater] Challenge @%s is in reveal. Sleeping until %s", twitterHandle, revealEndTime.Format(time.UnixDate))
			time.Sleep(time.Until(revealEndTime) + (2 * time.Minute))
		}

		if revealEndTime.Before(time.Now()) {
			updateStatus(application, errChan)
		}
	}

	// Fallthrough case is for applications that are whitelisted and have no
	// open challenges (we don't need to do anything for them)
}

func reveal(application *contracts.RegistryListing, errChan chan<- error) {
	votes, err := models.FindUnrevealedVotesFromPoll(application.ChallengeID.Int64())
	if err != nil {
		errChan <- err
		return
	}

	log.Printf("[updater] Revealing %d votes on poll %s.", len(votes), application.ChallengeID.String())
	for _, vote := range votes {
		account, err := models.FindAccountByID(vote.AccountID)
		if err != nil {
			errChan <- err
			continue
		} else if account == nil {
			errChan <- fmt.Errorf("Could not find account for ID %d", vote.AccountID)
			continue
		}

		log.Printf("\tRevealing %t vote by %s", vote.Vote, account.TwitterHandle)
		_, err = contracts.PLCRRevealVote(account.MultisigAddress.String, application.ChallengeID, vote.Vote, vote.Salt)
		if err != nil {
			errChan <- err
			continue
		}

		err = vote.MarkRevealed()
		if err != nil {
			errChan <- err
		}
	}
}

func updateStatus(application *contracts.RegistryListing, errChan chan<- error) {
	log.Printf("[updater] Attempting to updateStatus of listing 0x%x", application.ListingHash)
	// Refresh the listing, just in case there was a delay before calling this
	// function
	application, err := contracts.GetListingFromHash(application.ListingHash)
	if err != nil {
		errChan <- err
		return
	}

	// Reschedule if they have an ongoing challenge being waged against them
	// and it's not yet reveal time
	if application.ChallengeID.Cmp(big.NewInt(0)) != 0 {
		poll, err := contracts.GetPoll(application.ChallengeID)
		if err != nil {
			errChan <- err
			return
		}

		revealEndTime := time.Unix(poll.RevealEndDate.Int64(), 0)
		if revealEndTime.After(time.Now()) {
			go scheduleListing(application, errChan)
			return
		}
	}

	tx, err := contracts.UpdateStatus(application.ListingHash)
	if err != nil {
		errChan <- err
		return
	}
	log.Printf("[updater] Done! Updating tx: %s", tx.Hash().Hex())
}
