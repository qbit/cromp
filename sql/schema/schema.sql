CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
	user_id		BIGSERIAL	PRIMARY KEY,
	created_at	timestamp 	NOT NULL DEFAULT NOW(),
	updated_at	timestamp,
	first_name	text		NOT NULL,
	last_name	text		NOT NULL,
	username	text 		NOT NULL UNIQUE,
	hash		text 		NOT NULL,
	email		text		NOT NULL,
	token		UUID		NOT NULL default gen_random_uuid() UNIQUE,
	token_expires	timestamp	NOT NULL DEFAULT NOW() + INTERVAL '3 days'
);

CREATE TABLE entries (
	entry_id	UUID		NOT NULL default gen_random_uuid() PRIMARY KEY,
	user_id		int		NOT NULL REFERENCES users ON DELETE CASCADE,
	created_at	timestamp 	NOT NULL DEFAULT NOW(),
	updated_at	timestamp,
	title		text		NOT NULL DEFAULT '',
	body		text		NOT NULL DEFAULT ''
);

CREATE INDEX body_idx ON entries USING GIN (body gin_trgm_ops);

CREATE OR REPLACE FUNCTION hash(password text) RETURNS text AS $$
	SELECT crypt(password, gen_salt('bf', 10));
$$ LANGUAGE SQL;
