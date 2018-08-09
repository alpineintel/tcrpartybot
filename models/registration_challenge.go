package models

import (
	"time"
)

type RegistrationChallenge struct {
	ID                     int64 `db:"id"`
	AccountID              int64 `db:"account_id"`
	Account                *Account
	RegistrationQuestionID int64 `db:"registration_question_id"`
	RegistrationQuestion   *RegistrationQuestion
	SentAt                 *time.Time `db:"sent_at"`
	CompletedAt            *time.Time `db:"completed_at"`
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

func MarkRegistrationChallengeSent(challenge *RegistrationChallenge) error {
	db := GetDBSession()

	now := time.Now()
	_, err := db.Exec(`
		UPDATE registration_challenges
		SET sent_at=$1
		WHERE id=$2
	`, &now, challenge.ID)

	return err
}
