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
	token		text		NOT NULL default encode(digest(gen_random_uuid()::text || now(), 'sha256'), 'hex') UNIQUE,
	token_expires	timestamp	NOT NULL DEFAULT NOW() + INTERVAL '3 days'
);

CREATE TABLE entries (
	entry_id	UUID		NOT NULL default gen_random_uuid() PRIMARY KEY UNIQUE,
	user_id		BIGSERIAL	NOT NULL REFERENCES users ON DELETE CASCADE,
	created_at	timestamp 	NOT NULL DEFAULT NOW(),
	updated_at	timestamp,
	title		text		NOT NULL DEFAULT '',
	body		text		NOT NULL DEFAULT ''
);

CREATE INDEX body_trgm_idx ON entries USING gist (body gist_trgm_ops);

CREATE OR REPLACE FUNCTION hash(password text) RETURNS text AS $$
	SELECT crypt(password, gen_salt('bf', 10));
$$ LANGUAGE SQL;

-- CREATE OR REPLACE FUNCTION similar_entries(user_id bigserial, body text, OUT entry_id UUID, OUT similarity float, OUT headline text) AS $$
-- $$ LANGUAGE SQL;
