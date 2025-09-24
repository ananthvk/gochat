ALTER TABLE grp
ADD public_id bytea;

CREATE UNIQUE INDEX idx_public_id ON grp (public_id);