ALTER TABLE eth_events ADD COLUMN tx_hash VARCHAR;
ALTER TABLE eth_events ADD COLUMN tx_index INTEGER;
ALTER TABLE eth_events ADD COLUMN log_index INTEGER;
