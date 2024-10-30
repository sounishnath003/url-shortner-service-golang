-- queries.sql

-- name: GetUserByEmail
SELECT id, email, password 
FROM  users WHERE LOWER(email) = LOWER($1);

-- name: CreateNewUser
INSERT INTO users
(name, email, password) 
VALUES ($1,$2,$3);

-- name: CreateShortUrlQuery
INSERT INTO url_mappings
(original_url, short_url, expiration_at, user_id) 
VALUES ($1,$2,$3,$4);

-- name: GetShortUrlQuery
SELECT original_url, expiration_at FROM url_mappings 
WHERE short_url = $1 AND expiration_at > CURRENT_TIMESTAMP;

-- name: IncrUrlHitCountQuery
UPDATE url_mappings
SET hit_count = hit_count + 1
WHERE short_url = $1;

-- name: GetIncrementalIDQuery
select nextval('incr_id_generator_seq');

-- name: GetAllShortUrlAliasQuery
SELECT DISTINCT short_url
FROM url_mappings 
WHERE expiration_at > CURRENT_TIMESTAMP;