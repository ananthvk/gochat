ALTER TABLE message
ADD COLUMN sender_id BYTEA NOT NULL;

-- Note: There is no ON DELETE clause
-- This is to enable soft deleting of messages
ALTER TABLE message
ADD CONSTRAINT fk_message_sender_id
FOREIGN KEY (sender_id) REFERENCES usr(id);

CREATE INDEX IF NOT EXISTS idx_message_sender_id ON message(sender_id);