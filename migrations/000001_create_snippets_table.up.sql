CREATE TABLE snippets
(
    id      bigserial PRIMARY KEY,
    title   VARCHAR(100)                NOT NULL,
    content TEXT                        NOT NULL,
    created timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    expires timestamp(0) with time zone NOT NULL DEFAULT NOW() + INTERVAL '365 DAYS'
);

INSERT INTO snippets (title, content)
VALUES ('An old silent pond',
        E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō');

INSERT INTO snippets (title, content)
VALUES ('Over the wintry forest',
        E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki');

INSERT INTO snippets (title, content)
VALUES ('First autumn morning',
        E'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo');
