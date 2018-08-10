package models

import (
	"time"
)

const (
	REGISTRATION_CHALLENGE_COUNT = 3
)

type RegistrationChallenge struct {
	ID                     int64 `db:"id"`
	AccountID              int64 `db:"account_id"`
	Account                *Account
	RegistrationQuestionID int64                 `db:"registration_question_id"`
	RegistrationQuestion   *RegistrationQuestion `db:"registration_question"`
	SentAt                 *time.Time            `db:"sent_at"`
	CompletedAt            *time.Time            `db:"completed_at"`
}

func CreateRegistrationChallenge(account *Account, question *RegistrationQuestion) (*RegistrationChallenge, error) {
	db := GetDBSession()

	challenge := &RegistrationChallenge{
		AccountID:              account.ID,
		Account:                account,
		RegistrationQuestionID: question.ID,
		RegistrationQuestion:   question,
	}

	result := db.MustExec(`
		INSERT INTO registration_challenges (
			account_id,
			registration_question_id
		) VALUES($1, $2)
	`, account.ID, question.ID, time.Now())

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	challenge.ID = id

	return challenge, nil
}

func MarkRegistrationChallengeSent(challengeId int64) error {
	db := GetDBSession()

	now := time.Now()
	_, err := db.Exec(`
		UPDATE registration_challenges
		SET sent_at=$1
		WHERE id=$2
	`, &now, challengeId)

	return err
}

func FindUnsentChallenge(accountId int64) (*RegistrationChallenge, error) {
	db := GetDBSession()

	challenge := &RegistrationChallenge{}
	err := db.Get(&challenge, `
		SELECT
			registration_challenges.*,
			registration_questions.*
		FROM registration_challenges
		JOIN questions ON
			registration_challenges.question_id = registration_questions.id
		WHERE
			account_id = $1 AND
			sent_at IS NULL
		ORDER BY RANDOM()
		LIMIT 1
	`, accountId)

	return challenge, err
}

func FindIncompleteChallenge(accountId int64) (*RegistrationChallenge, error) {
	db := GetDBSession()

	challenge := &RegistrationChallenge{}
	err := db.Get(&challenge, `
		SELECT
			registration_challenges.*,
			registration_questions.*
		FROM registration_challenges
		JOIN questions ON
			registration_challenges.question_id = registration_questions.id
		WHERE
			account_id = $1 AND
			sent_at IS NOT NULL AND
			completed_at IS NULL
		LIMIT 1
	`, accountId)

	return challenge, err
}
