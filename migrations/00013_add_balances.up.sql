CREATE TABLE balances (
    id                         SERIAL PRIMARY KEY NOT NULL,
    account_id                 INTEGER NOT NULL,
    eth_event_id               INTEGER NOT NULL,
    plcr_balance               VARCHAR NOT NULL,
    wallet_balance             VARCHAR NOT NULL
);

ALTER TABLE ONLY public.balances
    ADD CONSTRAINT balances_accounts FOREIGN KEY (account_id)
    REFERENCES public.accounts(id);
ALTER TABLE ONLY public.balances
    ADD CONSTRAINT balances_eth_events FOREIGN KEY (eth_event_id)
    REFERENCES public.eth_events(id);
