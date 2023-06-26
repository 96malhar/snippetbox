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
VALUES ('Snippet 1 Title',
        'Snippet 1 content.',
        '2022-01-01 10:00:00',
        '2023-01-01 10:00:00');

INSERT INTO snippets (title, content, created, expires)
VALUES ('Snippet 2 Title',
        'Snippet 2 content.',
        '2022-02-01 10:00:00',
        '2023-02-01 10:00:00');

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

INSERT INTO users (name, email, hashed_password, created)
VALUES ('John',
        'john@example.com',
           -- Hello, World! as password
        '$2a$04$iQ07aWdTTLrEcem61mMEeuguBE994i.4qA5F90EhsPi9UQWzTBnyO',
        '2023-02-01 10:00:00');
