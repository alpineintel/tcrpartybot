CREATE TABLE eth_events (
    id                         SERIAL PRIMARY KEY NOT NULL,
    event_type                 VARCHAR NOT NULL,
    data                       JSON NOT NULL DEFAULT '{}',
    block_number               NUMERIC NOT NULL,
    created_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL
        DEFAULT (now() AT TIME ZONE 'utc')
);
