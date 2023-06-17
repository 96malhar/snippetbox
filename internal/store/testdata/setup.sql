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

INSERT INTO snippets (title, content, created, expires)
VALUES ('An old silent pond',
        E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
        '2022-01-01 10:00:00',
        '2023-01-01 10:00:00');

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
