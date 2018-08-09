package models

type RegistrationQuestion struct {
	ID       int64  `db:"id"`
	Question string `db:"question"`
	Answer   string `db:"answer"`
}

func FetchRandomRegistrationQuestions(count uint) ([]RegistrationQuestion, error) {
	db := GetDBSession()

	questions := []RegistrationQuestion{}
	err := db.Select(&questions, "SELECT * FROM registration_questions ORDER BY RANDOM() LIMIT ?", count)

	if err != nil {
		return nil, err
	}

	return questions, nil
}
