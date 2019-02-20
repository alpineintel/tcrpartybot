ALTER TABLE accounts DROP COLUMN multisig_address;

ALTER TABLE accounts ADD COLUMN eth_private_key VARCHAR(64);
ALTER TABLE accounts ADD COLUMN eth_address CHARACTER VARYING;
