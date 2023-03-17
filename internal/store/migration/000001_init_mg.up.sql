-- migrate -path ./internal/store/migration -database "postgresql://postgres@localhost:5432/snippetbox?sslmode=disable" -verbose up

--Create a `snippets` table.
CREATE TABLE snippets
(
    id      SERIAL PRIMARY KEY       NOT NULL,
    title   VARCHAR(100)             NOT NULL,
    content TEXT                     NOT NULL,
    created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    expires TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON snippets (created);

-- Add some dummy records (which we'll use in the next couple of chapters).
INSERT INTO snippets (title, content, created, expires)
VALUES ('An old silent pond',
        'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
        NOW() AT TIME ZONE 'UTC',
        (NOW() + INTERVAL '365 DAYS') AT TIME ZONE 'UTC');

INSERT INTO snippets (title, content, created, expires)
VALUES ('Over the wintry forest',
        'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
        NOW() AT TIME ZONE 'UTC',
        (NOW() + INTERVAL '365 DAYS') AT TIME ZONE 'UTC');

INSERT INTO snippets (title, content, created, expires)
VALUES ('First autumn morning',
        'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
        NOW() AT TIME ZONE 'UTC',
        (NOW() + INTERVAL '365 DAYS') AT TIME ZONE 'UTC');

CREATE USER web WITH password 'malhar123';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippets TO web;
GRANT USAGE, SELECT ON SEQUENCE snippets_id_seq TO web;
