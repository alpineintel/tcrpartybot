----
-- Accounts
--
CREATE TABLE IF NOT EXISTS accounts (
    id                               INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    twitter_handle                   VARCHAR(15) NOT NULL,
    eth_address                      VARCHAR(42) NOT NULL,
    eth_private_key                  VARCHAR(64) NOT NULL,
    passed_registration_challenge_at DATETIME,
    created_at             DATETIME  DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_handle ON accounts (twitter_handle);

----
-- Registration questions
--
CREATE TABLE IF NOT EXISTS registration_questions (
    id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT   NOT NULL
);

INSERT INTO registration_questions (question, answer) VALUES(
    "What is 2+2?",
    "4"
);

INSERT INTO registration_questions (question, answer) VALUES(
    "What is 4+1?",
    "5"
);

INSERT INTO registration_questions (question, answer) VALUES(
    "What is 10*5?",
    "50"
);

----
-- Registration challenges
--
CREATE TABLE IF NOT EXISTS registration_challenges (
    id                                    INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    account_id                            INTEGER NOT NULL,
    registration_question_id              INTEGER NOT NULL,
    sent_at                               DATETIME,
    completed_at                          DATETIME,
    FOREIGN KEY(registration_question_id) REFERENCES accounts(id),
    FOREIGN KEY(account_id)               REFERENCES accounts(id)
)
