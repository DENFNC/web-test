-- Миграция: инициализация схемы
CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(64) PRIMARY KEY,
    password_hash VARCHAR(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    token VARCHAR(128) PRIMARY KEY,
    login VARCHAR(64) NOT NULL REFERENCES users(login) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS documents (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    mime VARCHAR(64) NOT NULL,
    file BOOLEAN NOT NULL,
    public BOOLEAN NOT NULL,
    created TIMESTAMP NOT NULL,
    grant TEXT[] NOT NULL,
    owner VARCHAR(64) NOT NULL REFERENCES users(login) ON DELETE CASCADE,
    json JSONB
); 