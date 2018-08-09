package models

type ChallengeQuestion struct {
	ID       int64  `db:"id"`
	Question string `db:"question"`
	Answer   string `db:"answer"`
}

func FetchRandomChallengeQuestion() *ChallengeQuestion {
	db := GetDBSession()

	question := &ChallengeQuestion{}
	err := db.Get(question, "SELECT * FROM challenge_questions ORDER BY RANDOM() LIMIT 1")

	if err != nil {
		return nil
	}

	return question
}
