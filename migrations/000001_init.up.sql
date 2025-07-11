CREATE TABLE IF NOT EXISTS
    users (
        id UUID PRIMARY KEY,
        login TEXT UNIQUE NOT NULL CHECK (
            LENGTH(login) >= 8
            AND login ~ '^[a-zA-Z0-9]+$'
        ),
        password_hash TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    documents (
        id UUID PRIMARY KEY,
        file_name TEXT NOT NULL CHECK (LENGTH(file_name) > 0),
        mime_type TEXT NOT NULL DEFAULT 'application/octet-stream',
        has_file BOOLEAN NOT NULL DEFAULT FALSE,
        is_public BOOLEAN NOT NULL DEFAULT FALSE,
        owner_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS
    document_access (
        document_id UUID NOT NULL REFERENCES documents (id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        PRIMARY KEY (document_id, user_id)
    );

CREATE TABLE IF NOT EXISTS
    auth_tokens (
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        token TEXT PRIMARY KEY,
        is_revoked BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ DEFAULT NOW()
    );