CREATE TABLE snippets
(
    id      bigserial PRIMARY KEY,
    title   VARCHAR(100)                NOT NULL,
    content TEXT                        NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    expires timestamp(0) with time zone NOT NULL DEFAULT NOW() + INTERVAL '365 DAYS'
);

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
    id              bigserial PRIMARY KEY,
    name            varchar(255)                NOT NULL,
    email           varchar(255)                NOT NULL,
    hashed_password char(60)                    NOT NULL,
    created         timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

INSERT INTO users (name, email, hashed_password, created)
VALUES ('John',
        'john@example.com',
           -- Hello, World! as password
        '$2a$04$iQ07aWdTTLrEcem61mMEeuguBE994i.4qA5F90EhsPi9UQWzTBnyO',
        '2023-02-01 10:00:00');
