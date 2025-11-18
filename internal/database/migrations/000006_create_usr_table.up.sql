CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS usr (
    id BYTEA NOT NULL CHECK(length(id) = 16),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name TEXT NOT NULL,
    username TEXT NOT NULL,
    email CITEXT NOT NULL,
    password BYTEA NOT NULL,
    activated BOOL NOT NULL,

    CONSTRAINT Uk_usr_email UNIQUE (email),
    CONSTRAINT Uk_usr_username UNIQUE (username),
    CONSTRAINT Pk_usr PRIMARY KEY(id)
);