package twitter

import (
	"github.com/stevenleeg/go-twitter/twitter"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Stores all threads waiting to be unlocked
var requestQueue chan chan bool

// MonitorRatelimit will ensure that all twitter calls are executed at most
// once per two seconds
func MonitorRatelimit() {
	requestQueue = make(chan chan bool, 500)

	for {
		time.Sleep(2 * time.Second)
		request := <-requestQueue
		close(request)
	}
}

// awaitRatelimit will add the current thread's execution to a queue and will
// block until it is released by the ratelimiting thread
func awaitRatelimit() {
	await := make(chan bool)
	requestQueue <- await
	<-await
}

// GetIDFromHandle converts a twitter handle into an ID
func GetIDFromHandle(handle string) (int64, error) {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return 0, err
	}

	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: handle,
	})
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// FilterTweets will begin filtering tweets and outputting them to the returned
// channel
func FilterTweets(twitterIDs []string) (*twitter.Stream, error) {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return nil, err
	}

	params := &twitter.StreamFilterParams{
		Follow: twitterIDs,
	}

	stream, err := client.Streams.Filter(params)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// Retweet creates a new RT of the given tweet ID
func Retweet(handle string, tweetID int64) error {
	client, _, err := GetClientFromHandle(handle)
	if err != nil {
		return nil
	}

	awaitRatelimit()

	log.Printf("Retweeting from %s: %d", handle, tweetID)
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return nil
	}

	if _, _, err = client.Statuses.Retweet(tweetID, nil); err != nil {
		return err
	}

	return err
}

// SendTweet sends a new tweet from the given handle constant
func SendTweet(handle string, message string) error {
	client, _, err := GetClientFromHandle(handle)
	if err != nil {
		return nil
	}

	awaitRatelimit()

	log.Printf("Tweeting from %s: %s", handle, message)
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return nil
	}

	if _, _, err = client.Statuses.Update(message, nil); err != nil {
		return nil
	}

	return nil
}

// SendDM sends a direct message from the VIP party bot to the specified handle
func SendDM(recipientID int64, message string) error {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return err
	}

	awaitRatelimit()

	log.Printf("Sending DM to %d: %s", recipientID, message)
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return nil
	}

	_, _, err = client.DirectMessages.EventsNew(&twitter.DirectMessageEventsNewParams{
		Event: &twitter.DirectMessageEvent{
			Type: "message_create",
			Message: &twitter.DirectMessageEventMessage{
				Target: &twitter.DirectMessageTarget{
					RecipientID: strconv.FormatInt(recipientID, 10),
				},
				Data: &twitter.DirectMessageData{
					Text: message,
				},
			},
		},
	})

	if err != nil {
		log.Printf("Failed sending DM to %d: %s", recipientID, message)
		return err
	}

	return nil
}

// Follow will create a new friendship with the given user ID
func Follow(userID int64) error {
	log.Printf("Following Twitter user with ID %d", userID)
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return nil
	}

	awaitRatelimit()

	client, _, err := GetClientFromHandle(VIPBotHandle)

	follow := true
	params := &twitter.FriendshipCreateParams{
		UserID: userID,
		Follow: &follow,
	}

	if _, _, err = client.Friendships.Create(params); err != nil {
		return err
	}

	return nil
}

// IsFollower will return true if the given userID is a follower of the VIP bot
func IsFollower(userID int64) (bool, error) {
	if os.Getenv("SEND_TWITTER_INTERACTIONS") == "false" {
		return true, nil
	}

	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return false, err
	}

	params := &twitter.FriendshipShowParams{
		SourceID:         userID,
		TargetScreenName: os.Getenv("VIP_BOT_HANDLE"),
	}

	friendship, _, err := client.Friendships.Show(params)
	if err != nil {
		return false, err
	}

	return friendship.Source.Following, nil
}

// CreateWebhook creates a new webhook and subscribes it to the user, allowing
// us to receive notifications for new DMs. This should only be used on the
// TCRPartyVIP bot.
func CreateWebhook() (string, error) {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return "", err
	}

	webhookParams := &twitter.AccountActivityRegisterWebhookParams{
		EnvName: os.Getenv("TWITTER_ENV"),
		URL:     os.Getenv("BASE_URL") + "/webhooks/twitter",
	}
	webhook, _, err := client.AccountActivity.RegisterWebhook(webhookParams)

	if err != nil {
		return "", err
	}

	return webhook.ID, nil
}

// CreateSubscription subscribes the current webhook to the given user
func CreateSubscription() error {
	client, _, err := GetClientFromHandle(VIPBotHandle)
	if err != nil {
		return err
	}

	subParams := &twitter.AccountActivityCreateSubscriptionParams{
		EnvName: os.Getenv("TWITTER_ENV"),
	}
	_, err = client.AccountActivity.CreateSubscription(subParams)

	return err
}

// GetList fetches the current members of the TCR list on Twitter
func GetList() (*twitter.Members, error) {
	client, _, err := GetClientFromHandle(PartyBotHandle)
	if err != nil {
		return nil, err
	}

	members, _, err := client.Lists.Members(&twitter.ListsMembersParams{
		Slug:            os.Getenv("TWITTER_LIST_SLUG"),
		OwnerScreenName: os.Getenv("PARTY_BOT_HANDLE"),
	})
	if err != nil {
		return nil, err
	}

	return members, nil
}

// AddHandlesToList takes a slice of Twitter handles and adds them to the TCR
// list on Twitter
func AddHandlesToList(handles []string) error {
	client, _, err := GetClientFromHandle(PartyBotHandle)
	if err != nil {
		return err
	}

	screenNames := strings.Join(handles, ",")
	_, err = client.Lists.MembersCreateAll(&twitter.ListsMembersCreateAllParams{
		Slug:            os.Getenv("TWITTER_LIST_SLUG"),
		OwnerScreenName: os.Getenv("PARTY_BOT_HANDLE"),
		ScreenName:      screenNames,
	})
	return err
}

// RemoveHandlesToList takes a slice of Twitter handles and removes them to the
// TCR list on Twitter
func RemoveHandlesToList(handles []string) error {
	client, _, err := GetClientFromHandle(PartyBotHandle)
	if err != nil {
		return err
	}

	screenNames := strings.Join(handles, ",")
	_, err = client.Lists.MembersDestroyAll(&twitter.ListsMembersDestroyAllParams{
		Slug:            os.Getenv("TWITTER_LIST_SLUG"),
		OwnerScreenName: os.Getenv("PARTY_BOT_HANDLE"),
		ScreenName:      screenNames,
	})
	return err
}

// GetListTweets retrieves a slice of tweets created since the last tweet ID
// provided
func GetListTweets(sinceID int64) ([]twitter.Tweet, error) {
	client, _, err := GetClientFromHandle(PartyBotHandle)
	if err != nil {
		return nil, err
	}

	fls := false
	tweets, _, err := client.Lists.Statuses(&twitter.ListsStatusesParams{
		Slug:            os.Getenv("TWITTER_LIST_SLUG"),
		OwnerScreenName: os.Getenv("PARTY_BOT_HANDLE"),
		SinceID:         sinceID,
		IncludeRetweets: &fls,
	})

	return tweets, err
}
