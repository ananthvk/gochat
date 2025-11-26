CREATE TABLE IF NOT EXISTS grp_membership (
    grp_id   BYTEA NOT NULL,
    usr_id   BYTEA NOT NULL,
    role     TEXT NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT Pk_grp_membership PRIMARY KEY (grp_id, usr_id),
    CONSTRAINT Fk_grp_membership_grp FOREIGN KEY (grp_id) REFERENCES grp(id) ON DELETE CASCADE,
    CONSTRAINT Fk_grp_membership_usr FOREIGN KEY (usr_id) REFERENCES usr(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_grp_membership_usr_id ON grp_membership(usr_id);

CREATE INDEX IF NOT EXISTS idx_grp_membership_grp_id ON grp_membership(grp_id);