package models

import (
	"time"
)

// RegistrationQuestion is a struct representing a row in the
// registration_questions table
type RegistrationQuestion struct {
	ID        int64      `db:"id"`
	Question  string     `db:"question"`
	Answer    string     `db:"answer"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// FetchRandomRegistrationQuestions will return `count` number of registration
// questions from the database
func FetchRandomRegistrationQuestions(count uint) ([]RegistrationQuestion, error) {
	db := GetDBSession()

	questions := []RegistrationQuestion{}
	err := db.Select(&questions, "SELECT * FROM registration_questions WHERE deleted_at IS NULL ORDER BY RANDOM() LIMIT $1", count)

	if err != nil {
		return nil, err
	}

	return questions, nil
}
