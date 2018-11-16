ALTER TABLE accounts ADD COLUMN multisig_factory_identifier BIGINT;
ALTER TABLE accounts DROP eth_private_key;
ALTER TABLE accounts DROP eth_address;
