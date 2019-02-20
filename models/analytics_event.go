package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

// AnalyticsEventType represents the various event types in our analytics
// table. This reflects the analyitics_event_type enum in postgres.
type AnalyticsEventType string

const (
	DMAnalyticsEvent      AnalyticsEventType = "dm"
	MentionAnalyticsEvent AnalyticsEventType = "mention"
)

// AnalyticsEvent represents any form of event that we may want to keep track
// of for product metrics
type AnalyticsEvent struct {
	ID          int64              `db:"id"`
	Key         AnalyticsEventType `db:"key"`
	AccountID   sql.NullInt64      `db:"account_id"`
	Additionals types.JSONText     `db:"additionals"`
	CreatedAt   *time.Time         `db:"created_at"`
}

type dmEventData struct {
	TwitterUserID int64  `json:"twitter_user_id"`
	Msg           string `json:"msg"`
}

type mentionEventData struct {
	TwitterUserID int64  `json:"twitter_user_id"`
	Msg           string `json:"msg"`
	TweetID       string `json:"tweet_id"`
}

type ethEventData struct {
	EventName string      `json:"event_name"`
	EventData interface{} `json:"event_data"`
}

func createAnalyticsEvent(event *AnalyticsEvent) error {
	db := GetDBSession()

	var id int64
	err := db.QueryRow(`
		INSERT INTO analytics_events (
			key,
			account_id,
			additionals
		) VALUES($1, $2, $3)
		RETURNING id
	`, event.Key, event.AccountID, event.Additionals).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

func generateAccountIDInt(twitterID int64) (*sql.NullInt64, error) {
	// Does this handle have an account?
	account, err := FindAccountByTwitterID(twitterID)
	if err != nil {
		return nil, err
	}

	accountIDInt := sql.NullInt64{Valid: false}
	if account != nil {
		accountIDInt.Int64 = account.ID
		accountIDInt.Valid = true
	}

	return &accountIDInt, nil
}

// CreateDMAnalyticsEvent creates a new event given an account ID (optional,
// set to 0 if null) and the direct message's text
func CreateDMAnalyticsEvent(twitterID int64, message string) error {
	additionals := dmEventData{twitterID, message}
	bytes, err := json.Marshal(additionals)
	if err != nil {
		return err
	}

	accountIDInt, err := generateAccountIDInt(twitterID)
	if err != nil {
		return err
	}

	event := &AnalyticsEvent{
		Key:         DMAnalyticsEvent,
		AccountID:   *accountIDInt,
		Additionals: types.JSONText(string(bytes)),
	}

	return createAnalyticsEvent(event)
}

// CreateMentionAnalyticsEvent creates a new event given an account ID
// (optional, set to 0 if null), the tweet's text, and its ID.
func CreateMentionAnalyticsEvent(twitterID int64, tweetID string, message string) error {
	additionals := mentionEventData{twitterID, message, tweetID}
	bytes, err := json.Marshal(additionals)
	if err != nil {
		return err
	}

	accountIDInt, err := generateAccountIDInt(twitterID)
	if err != nil {
		return err
	}

	event := &AnalyticsEvent{
		Key:         MentionAnalyticsEvent,
		AccountID:   *accountIDInt,
		Additionals: types.JSONText(string(bytes)),
	}

	return createAnalyticsEvent(event)
}
