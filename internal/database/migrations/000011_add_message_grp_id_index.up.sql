-- Makes it quicker to access the last message of a group
CREATE INDEX idx_message_grp_id_id_desc ON message(grp_id, id DESC);