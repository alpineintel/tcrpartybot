CREATE TABLE public.votes (
    poll_id       INTEGER NOT NULL,
    account_id    INTEGER NOT NULL,
    salt          INTEGER NOT NULL,
    vote          BOOLEAN NOT NULL,
    revealed_at   TIMESTAMP WITHOUT TIME ZONE,
    created_at    TIMESTAMP WITHOUT TIME ZONE NOT NULL 
        DEFAULT (now() AT TIME ZONE 'utc'),
    PRIMARY KEY (poll_id, account_id)
);

ALTER TABLE ONLY public.votes
    ADD CONSTRAINT vote_account FOREIGN KEY (account_id)
    REFERENCES public.accounts(id);
