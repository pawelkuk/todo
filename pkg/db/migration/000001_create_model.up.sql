CREATE TABLE IF NOT EXISTS task(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    description TEXT,
    created_at INTEGER,
    updated_at INTEGER,
    due_date INTEGER,
    completed INTEGER
);
CREATE TABLE IF NOT EXISTS periodic_task(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    description TEXT,
    created_at INTEGER,
    updated_at INTEGER,
    schedule TEXT
);
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE,
    password TEXT,
    authorization_token TEXT
);
CREATE TABLE session (
    token TEXT PRIMARY KEY,
    expiry REAL NOT NULL,
    user_id TEXT,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
);
CREATE INDEX sessions_expiry_idx ON session(expiry);