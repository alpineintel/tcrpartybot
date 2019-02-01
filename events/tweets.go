package events

import (
	"log"

	goTwitter "github.com/stevenleeg/go-twitter/twitter"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

// ListenAndRetweet listens for all tweets by users who are on the TCR and
// retweets onto @TCRParty's timeline. It will also listen for updates on the
// ethEvents channel and reset the twitter filter upon hearing of an
// addition/removal of a listing on the TCR
func ListenAndRetweet(ethEvents <-chan *ETHEvent, errChan chan<- error) {
	var stream *goTwitter.Stream

	refreshListings := func() {
		listings, err := contracts.GetWhitelistedListings()
		if err != nil {
			errChan <- err
			return
		}

		// Convert listings into twitter handles
		var twitterHandles []string
		handleMap := map[string]bool{}
		for _, listing := range listings {
			handle, err := contracts.GetListingDataFromHash(listing.ListingHash)
			if err != nil {
				errChan <- err
				continue
			}

			handleMap[handle] = true
			twitterHandles = append(twitterHandles, handle)
		}

		twitterStream, tweetChan, err := twitter.FilterTweets(twitterHandles)
		if err != nil {
			errChan <- err
			return
		}

		stream = twitterStream
		for tweet := range tweetChan {
			// Make sure we've received a tweet from a user we're tracking
			if !handleMap[tweet.User.ScreenName] {
				continue
			}

			// Don't retweet replies or retweets
			if tweet.InReplyToStatusIDStr != "" || tweet.RetweetedStatus != nil {
				continue
			}

			log.Printf("New tweet from @%s", tweet.User.ScreenName)
			err = twitter.Retweet(twitter.PartyBotHandle, tweet.ID)
			if err != nil {
				errChan <- err
				continue
			}
		}
		log.Println("Killing filter stream")
	}

	go refreshListings()

	for {
		event := <-ethEvents
		// Listen for whitelisted/removal events
		name := event.EventType
		if name != ETHEventTCRApplicationWhitelisted && name != ETHEventTCRListingRemoved {
			continue
		}

		log.Println("TCR update detected, refreshing list for retweets...")
		// Stop the existing stream and start up a new one
		if stream != nil {
			stream.Stop()
		}

		go refreshListings()
	}
}
