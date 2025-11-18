ALTER TABLE grp DROP CONSTRAINT IF EXISTS Pk_grp;

ALTER TABLE grp DROP COLUMN IF EXISTS id;

ALTER TABLE grp RENAME COLUMN public_id TO id;

ALTER TABLE grp ADD CONSTRAINT pk_grp PRIMARY KEY (id);

-- Redundant since primary key constraint creates an index
DROP INDEX IF EXISTS idx_public_id;