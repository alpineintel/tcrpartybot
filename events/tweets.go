package events

import (
	"log"
	"strconv"

	goTwitter "github.com/stevenleeg/go-twitter/twitter"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

var listingHashTwitterID = map[[32]byte]string{}
var idHandle = map[string]string{}

func fetchListings() ([]string, error) {
	listings, err := contracts.GetWhitelistedListings()
	if err != nil {
		return nil, err
	}

	// Convert listings into twitter IDs
	var twitterIDs []string
	log.Println("[Tweets] Filtering tweets by:")
	for _, listing := range listings {
		if id := listingHashTwitterID[listing.ListingHash]; id != "" {
			twitterIDs = append(twitterIDs, id)
			log.Printf("\t%s (cache)", idHandle[id])
			continue
		}

		handle, err := contracts.GetListingDataFromHash(listing.ListingHash)
		if err != nil {
			return nil, err
		}

		log.Printf("\t%s", handle)

		// Convert the handle to an ID
		id, err := twitter.GetIDFromHandle(handle)
		if err != nil {
			return nil, err
		}

		strID := strconv.FormatInt(id, 10)
		listingHashTwitterID[listing.ListingHash] = strID
		idHandle[strID] = handle
		twitterIDs = append(twitterIDs, strID)
	}

	return twitterIDs, nil
}

func handleTweet(tweet *goTwitter.Tweet) error {
	// Make sure we've received a tweet from a user we're tracking
	listingHash := contracts.GetListingHash(tweet.User.ScreenName)
	if listingHashTwitterID[listingHash] == "" {
		return nil
	}

	// Don't retweet replies or retweets
	if tweet.InReplyToStatusIDStr != "" || tweet.RetweetedStatus != nil {
		return nil
	}

	log.Printf("New tweet from @%s", tweet.User.ScreenName)
	err := twitter.Retweet(twitter.PartyBotHandle, tweet.ID)
	if err != nil {
		return err
	}

	return nil
}

// ListenAndRetweet listens for all tweets by users who are on the TCR and
// retweets onto @TCRParty's timeline. It will also listen for updates on the
// ethEvents channel and reset the twitter filter upon hearing of an
// addition/removal of a listing on the TCR
func ListenAndRetweet(ethEvents <-chan *ETHEvent, errChan chan<- error) {
	twitterIDs, err := fetchListings()
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	stream, err := twitter.FilterTweets(twitterIDs)
	if err != nil {
		errChan <- errors.Wrap(err)
		return
	}

	reset := make(chan bool)

	// Start listening for tweets and retweeting
	go func() {
		for {
			select {
			case msg := <-stream.Messages:
				tweet, ok := msg.(*goTwitter.Tweet)
				if !ok {
					continue
				}

				handleTweet(tweet)
			case <-reset:
				return
			}
		}
	}()

	// Listen for TCR changes that require us to reset the Twitter stream with
	// a new filter
	for {
		event := <-ethEvents
		name := event.EventType
		if name == ETHEventTCRApplicationWhitelisted || name == ETHEventTCRListingRemoved {
			break
		}
	}

	log.Println("TCR update detected, refreshing list for retweets...")
	reset <- true
	stream.Stop()
	go ListenAndRetweet(ethEvents, errChan)
}
