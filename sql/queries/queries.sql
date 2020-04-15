-- name: CreateEntry :one
INSERT INTO entries (
	  entry_id, user_id, title, body
) VALUES (
  $1, $2, $3, $4
)
RETURNING entry_id, created_at, to_tsvector(body);

-- name: UpdateEntry :execrows
UPDATE entries SET
	title = $3,
	body = $4
WHERE entry_id = $1 and
user_id = $2;

-- name: GetEntry :one
SELECT * FROM entries
WHERE entry_id = $1 and user_id = $2
LIMIT 1;

-- name: GetEntries :many
SELECT * FROM entries
WHERE user_id = $1;

-- name: DeleteEntry :execrows
DELETE FROM entries
WHERE entry_id = $1;

-- name: SimilarEntries :many
SELECT entry_id, similarity(body, $2) as similarity,
	ts_headline('english', body, q, 'StartSel = <b>, StopSel = </b>') as headline,
	title from entries,
	to_tsquery($2) q
WHERE user_id = $1 and
	similarity(body, $2) > 0.0
	order by similarity DESC
	LIMIT 10;

-- name: CreateUser :one
INSERT INTO users (
	  first_name, last_name, username, email, hash
) VALUES (
  $1, $2, $3, $4, hash($5)
)
RETURNING user_id, username, token, token_expires;

-- name: GetUser :one
SELECT * FROM users
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByToken :one
SELECT * FROM users
WHERE token = $1 LIMIT 1;

-- name: AuthUser :one
UPDATE users
SET
token = DEFAULT,
token_expires = DEFAULT
WHERE
username = $2 and
(hash = crypt($1, hash)) = true
RETURNING user_id, created_at, first_name, last_name, username, email, token, token_expires, true as authed;

-- name: DeleteUser :exec
DELETE FROM users CASCADE
WHERE user_id = $1;

-- name: ValidToken :one
SELECT now() < token_created FROM users
WHERE token = $1 LIMIT 1;

-- name: EntriesByToken :many
SELECT * FROM entries
WHERE user_id = (SELECT user_id FROM users WHERE token = $1 limit 1);
