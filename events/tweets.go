package events

import (
	"log"
	"strconv"
	"time"

	goTwitter "github.com/stevenleeg/go-twitter/twitter"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

func handleTweet(tweet goTwitter.Tweet) error {
	log.Printf("New tweet from @%s: %s", tweet.User.ScreenName, tweet.Text)

	// Don't retweet replies or retweets
	if tweet.InReplyToStatusIDStr != "" || tweet.RetweetedStatus != nil {
		log.Printf("\tSkipping reply tweet")
		return nil
	}

	err := twitter.Retweet(twitter.PartyBotHandle, tweet.ID)
	if err != nil {
		return err
	}

	return nil
}

func pollList(errChan chan<- error) {
	for {
		time.Sleep(10 * time.Second)

		// Get the last ID we've synced
		sinceIDStr, err := models.GetKey(models.LatestSyncedTweetKey)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}

		sinceID, _ := strconv.ParseInt(sinceIDStr, 10, 64)

		// Get any new tweets
		tweets, err := twitter.GetListTweets(sinceID)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}

		// Retweet all of these new ones
		for _, tweet := range tweets {
			err = handleTweet(tweet)
			if err != nil {
				errChan <- errors.Wrap(err)
				continue
			}
		}

		if len(tweets) > 0 {
			err = models.SetKey(models.LatestSyncedTweetKey, tweets[0].IDStr)
			if err != nil {
				errChan <- errors.Wrap(err)
				continue
			}
		}
	}
}

// ListenAndRetweet listens for all tweets by users who are on the TCR and
// retweets onto @TCRParty's timeline. It will also poll for updates on the
// eth_events table and reset the twitter filter upon hearing of an
// addition/removal of a listing on the TCR
func ListenAndRetweet(errChan chan<- error) {
	// Sync up the list just in case anything has happened since we've last
	// booted the bot
	twitter.SyncList()

	// Start polling the list and retweeting
	go pollList(errChan)

	// Listen for TCR changes that require us to reset the Twitter stream with
	// a new filter
	for {
		time.Sleep(10 * time.Second)

		// Fetch the last synced event ID
		latestEventIDStr, err := models.GetKey(models.LatestListSyncedEventKey)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}

		latestEventID := int64(0)
		if latestEventIDStr != "" {
			intID, err := strconv.ParseInt(latestEventIDStr, 10, 64)
			if err != nil {
				errChan <- errors.Wrap(err)
				continue
			}

			latestEventID = intID
		}

		needsUpdate, newestID, err := models.TCRHasUpdatedSinceEventID(latestEventID)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		} else if needsUpdate {
			twitter.SyncList()
		}

		// Set the last synced event ID
		if err := models.SetKey(models.LatestListSyncedEventKey, strconv.FormatInt(newestID, 10)); err != nil {
			errChan <- errors.Wrap(err)
			continue
		}
	}
}
