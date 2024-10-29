DROP TABLE IF EXISTS users_url_mappings;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS url_mappings;

-- users
CREATE TABLE IF NOT EXISTS users (
	id INT PRIMARY KEY GENERATED BY DEFAULT AS Identity,
	name VARCHAR(20) NOT NULL,
	email VARCHAR(20) NOT NULL UNIQUE,
	password VARCHAR(20) NOT NULL,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS users_email_idx ON users (email);

-- url_mappings
CREATE TABLE IF NOT EXISTS url_mappings (
	id INT PRIMARY KEY GENERATED BY DEFAULT AS Identity,
    user_id INT NOT NULL,
	original_url VARCHAR(255) NOT NULL,
	short_url VARCHAR(20) NOT NULL UNIQUE,
	expiration_at TIMESTAMP WITH TIME ZONE,
    createdAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS url_mappings_short_url_idx ON url_mappings (short_url);
CREATE INDEX IF NOT EXISTS url_mappings_user_id_idx ON url_mappings (user_id);

-- user_url_mappings
CREATE TABLE IF NOT EXISTS users_url_mappings (
    UrlID INT NOT NULL,
    UserID INT NOT NULL,
    PRIMARY KEY (UrlID, UserID),
    FOREIGN KEY (UrlID) REFERENCES url_mappings(id) ON DELETE CASCADE,
    FOREIGN KEY (UserID) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX users_url_mappings_userid_urlid_idx ON users_url_mappings(UserID, UrlID);

-- Add some data
INSERT INTO users (name, email, password) VALUES ('Sounish', 'sounish@example.com', 'password');

INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ('https://www.google.com', 'googl', '2024-01-01 00:00:00', 1);
INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ('https://www.facebook.com', 'faceb', '2024-01-01 00:00:00', 1);
INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ('https://www.twitter.com', 'twitt', '2024-01-01 00:00:00', 1);
INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ('https://www.instagram.com', 'insta', '2024-01-01 00:00:00', 1);
INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ('https://www.linkedin.com', 'linkd', '2024-01-01 00:00:00', 1);
INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ('https://www.youtube.com', 'youtb', '2024-01-01 00:00:00', 1);
INSERT INTO url_mappings (original_url, short_url, expiration_at, user_id) VALUES ('https://www.amazon.com', 'amazn', '2024-01-01 00:00:00', 1);

INSERT INTO users_url_mappings (UrlID, UserID) VALUES (1, 1);
INSERT INTO users_url_mappings (UrlID, UserID) VALUES (2, 1);
INSERT INTO users_url_mappings (UrlID, UserID) VALUES (3, 1);
INSERT INTO users_url_mappings (UrlID, UserID) VALUES (4, 1);
INSERT INTO users_url_mappings (UrlID, UserID) VALUES (5, 1);
INSERT INTO users_url_mappings (UrlID, UserID) VALUES (6, 1);
INSERT INTO users_url_mappings (UrlID, UserID) VALUES (7, 1);


-- Generate a incremental id generator
CREATE UNLOGGED SEQUENCE IF NOT EXISTS public.incr_id_generator_seq
    INCREMENT 1
    START 1000000
    MAXVALUE 999999999;

ALTER SEQUENCE public.incr_id_generator_seq
    OWNER TO root;

COMMENT ON SEQUENCE public.incr_id_generator_seq
    IS 'to generate the incremental id for short urls';