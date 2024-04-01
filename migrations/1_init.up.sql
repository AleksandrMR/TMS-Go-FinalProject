CREATE TABLE IF NOT EXISTS hashService
(
    hash    TEXT NOT NULL UNIQUE,
    payload TEXT
);
CREATE INDEX IF NOT EXISTS idx_hash ON hashService (hash);
