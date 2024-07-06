CREATE TABLE IF NOT EXISTS event (
    id text NOT NULL,
    pubkey text NOT NULL,
    created_at integer NOT NULL,
    kind integer NOT NULL,
    tags jsonb NOT NULL,
    content text NOT NULL,
    sig text NOT NULL
);