CREATE TYPE analytics_event_type AS ENUM (
    'dm',
    'mention'
);

CREATE TABLE public.analytics_events (
    id           SERIAL PRIMARY KEY NOT NULL,
    key          analytics_event_type NOT NULL,
    account_id   INTEGER,
    additionals  JSON NOT NULL DEFAULT '{}',
    created_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL
        DEFAULT (now() AT TIME ZONE 'utc')
);

ALTER TABLE ONLY public.analytics_events
    ADD CONSTRAINT analytics_events_account FOREIGN KEY (account_id)
    REFERENCES public.accounts(id);
