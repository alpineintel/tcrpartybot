CREATE TABLE IF NOT EXISTS accounts (
    id              INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    twitter_handle  VARCHAR(15) NOT NULL,
    eth_address     VARCHAR(42) NOT NULL,
    eth_private_key VARCHAR(64) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_handle ON accounts (twitter_handle);
