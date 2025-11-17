CREATE TABLE IF NOT EXISTS message (
    -- ID of the message
    id BYTEA NOT NULL CHECK(length(id) = 16),
    -- Type of the message - text, image, etc
    type TEXT NOT NULL,
    -- The group in which this message was sent
    grp_id BYTEA NOT NULL CHECK(length(grp_id) = 16),
    -- Time at which the message was created
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- Content of the message
    content TEXT NOT NULL,
    
    -- Many to one relationship, i.e. each group can have multiple messages, and one message can be part of one group
    CONSTRAINT Fk_message_grp FOREIGN KEY (grp_id) REFERENCES grp(id),
    CONSTRAINT Pk_message PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS idx_message_grp_id ON message(grp_id);