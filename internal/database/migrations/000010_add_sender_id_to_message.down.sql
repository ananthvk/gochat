DROP INDEX IF EXISTS idx_message_sender_id;

ALTER TABLE message
DROP CONSTRAINT IF EXISTS fk_message_sender_id;

ALTER TABLE message
DROP COLUMN IF EXISTS sender_id;