-- queries.sql

-- name: GetUserByEmail
SELECT email, password 
FROM  users WHERE LOWER(email) = LOWER($1);

-- name: CreateShortUrlQuery
INSERT INTO url_mappings
(original_url, short_url, expiration_at, user_id) 
VALUES ($1,$2,$3,$4);

-- name: GetShortUrlQuery
SELECT short_url, expiration_at FROM url_mappings WHERE short_url = $1;

-- name: GetIncrementalIDQuery
select nextval('incr_id_generator_seq');
