CREATE TABLE public.faucet_hits (
    id          SERIAL PRIMARY KEY NOT NULL,
    amount      VARCHAR NOT NULL,
    account_id  INTEGER NOT NULL,
    created_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL
        DEFAULT (now() AT TIME ZONE 'utc')
);

ALTER TABLE ONLY public.faucet_hits
    ADD CONSTRAINT faucet_hits_accounts FOREIGN KEY (account_id)
    REFERENCES public.accounts(id);
