-- migrate -path ./internal/store/migration -database "postgresql://postgres@localhost:5432/snippetbox?sslmode=disable" -verbose up

--Create a `snippets` table.
CREATE TABLE snippets
(
    id      INT Generated Always As IDENTITY PRIMARY KEY NOT NULL,
    title   VARCHAR(100)                                 NOT NULL,
    content TEXT                                         NOT NULL,
    created TIMESTAMP WITH TIME ZONE                     NOT NULL,
    expires TIMESTAMP WITH TIME ZONE                     NOT NULL
);
-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON snippets (created);

--Create a `users` table.
CREATE TABLE users
(
    id              INT Generated Always As IDENTITY PRIMARY KEY NOT NULL,
    name            VARCHAR(255)                                 NOT NULL,
    email           VARCHAR(255)                                 NOT NULL,
    hashed_password CHAR(60)                                     NOT NULL,
    created         TIMESTAMP WITH TIME ZONE                     NOT NULL,
    CONSTRAINT users_uc_email UNIQUE (email)
);

-- Create a `sessions` table.
CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BYTEA                    NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Add an index on the expiry column.
CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO snippets (title, content, created, expires)
VALUES ('An old silent pond',
        E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
        NOW(),
        NOW() + INTERVAL '365 DAYS');

INSERT INTO snippets (title, content, created, expires)
VALUES ('Over the wintry forest',
        E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
        NOW(),
        NOW() + INTERVAL '365 DAYS');

INSERT INTO snippets (title, content, created, expires)
VALUES ('First autumn morning',
        E'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
        NOW(),
        NOW() + INTERVAL '365 DAYS');

CREATE
    USER web WITH password 'malhar123';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippets TO web;
GRANT SELECT, INSERT, UPDATE, DELETE ON sessions TO web;
GRANT SELECT, INSERT, UPDATE, DELETE ON users TO web;
