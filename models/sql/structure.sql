CREATE TABLE IF NOT EXISTS accounts (
    id              INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    twitter_handle  VARCHAR(15) NOT NULL,
    eth_address     VARCHAR(42) NOT NULL,
    eth_private_key VARCHAR(64) NOT NULL,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_handle ON accounts (twitter_handle);

CREATE TABLE IF NOT EXISTS challenge_questions (
    id            INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT   NOT NULL
);

INSERT INTO challenge_questions (question, answer) VALUES(
    "What is 2+2?",
    "4"
);
