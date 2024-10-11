PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA auto_vacuum = FULL;
PRAGMA encoding = "UTF-8";
PRAGMA foreign_keys = ON;
-- Event
CREATE TABLE IF NOT EXISTS event (
    id TEXT PRIMARY KEY,
    pubkey TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    kind INTEGER NOT NULL,
    -- TODO: remove
    tags jsonb NOT NULL,
    content TEXT NOT NULL,
    sig TEXT NOT NULL,
    -- for Search
    profname TEXT NOT NULL COLLATE NOCASE
);
CREATE INDEX IF NOT EXISTS pubkeyidx ON event(pubkey);
CREATE INDEX IF NOT EXISTS createdidx ON event(created_at DESC);
CREATE INDEX IF NOT EXISTS kindidx ON event(kind);
CREATE INDEX IF NOT EXISTS kindcreatedidx ON event(kind, created_at DESC);
CREATE INDEX IF NOT EXISTS expiresidx ON event(expires_at);
CREATE INDEX IF NOT EXISTS profnameidx ON event(profname COLLATE NOCASE);
CREATE INDEX IF NOT EXISTS profcreatedidx ON event(profname COLLATE NOCASE, created_at DESC);
-- Tag for index
CREATE TABLE IF NOT EXISTS tag (
    event_id TEXT NOT NULL,
    name TEXT,
    first TEXT,
    FOREIGN KEY(event_id) REFERENCES event(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS tageventidx ON tag(event_id);
CREATE INDEX IF NOT EXISTS tagnameidx ON tag(name);
CREATE INDEX IF NOT EXISTS tagfirstidx ON tag(first);
CREATE INDEX IF NOT EXISTS tagnamefirstidx ON tag(name, first);