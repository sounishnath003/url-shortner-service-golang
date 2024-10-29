-- queries.sql

-- name: CreateShortUrlQuery
-- param: some
INSERT INTO url_mappings 
(original_url, short_url, expiration_at, user_id)
VALUES (?,?,?,?);

-- name: GetShortUrlQuery
-- param: some
SELECT short_url, expiration_at FROM url_mappings WHERE short_url = ?;

-- name: GetIncrementalIDQuery
select nextval('incr_id_generator_seq');
