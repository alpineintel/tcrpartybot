CREATE TABLE registry_listings (
    id                      VARCHAR NOT NULL PRIMARY KEY,
    listing_hash            VARCHAR NOT NULL,
    twitter_handle          VARCHAR NOT NULL,
    deposit                 VARCHAR NOT NULL,
    application_created_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    application_ended_at     TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    whitelisted_at          TIMESTAMP WITHOUT TIME ZONE,
    removed_at              TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE registry_challenges (
    poll_id                 INTEGER NOT NULL PRIMARY KEY,
    listing_hash            VARCHAR NOT NULL,
    listing_id              VARCHAR NOT NULL,
    challenger_account_id   INTEGER,
    challenger_address      VARCHAR NOT NULL,
    created_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    commit_ends_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    reveal_ends_at          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    succeeded_at            TIMESTAMP WITHOUT TIME ZONE,
    failed_at               TIMESTAMP WITHOUT TIME ZONE
);

ALTER TABLE ONLY public.registry_challenges
    ADD CONSTRAINT registry_challenge_account FOREIGN KEY (challenger_account_id)
    REFERENCES public.accounts(id);

ALTER TABLE ONLY public.registry_challenges
    ADD CONSTRAINT registry_challenge_registry_listing FOREIGN KEY (listing_id)
    REFERENCES public.registry_listings(id);
