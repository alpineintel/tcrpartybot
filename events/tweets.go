package events

import (
	"log"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

// ListenAndRetweet listens for all tweets by users who are on the TCR and
// retweets onto @TCRParty's timeline
func ListenAndRetweet(errChan chan<- error) {
	listings, err := contracts.GetWhitelistedListings()
	if err != nil {
		errChan <- err
		return
	}

	// Convert listings into twitter handles
	var twitterHandles []string
	for _, listing := range listings {
		handle, err := contracts.GetListingDataFromHash(listing.ListingHash)
		if err != nil {
			errChan <- err
			continue
		}

		twitterHandles = append(twitterHandles, handle)
	}

	_, tweetChan, err := twitter.FilterTweets(0, twitterHandles)
	if err != nil {
		errChan <- err
		return
	}

	for tweet := range tweetChan {
		log.Printf("New tweet from @%s", tweet.User.ScreenName)

		err = twitter.Retweet(twitter.PartyBotHandle, tweet.ID)
		if err != nil {
			errChan <- err
			continue
		}
	}
}
