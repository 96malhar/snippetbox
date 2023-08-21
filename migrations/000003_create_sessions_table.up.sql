CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BYTEA                    NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL
);
