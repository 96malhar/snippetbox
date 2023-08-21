CREATE TABLE users
(
    id              bigserial PRIMARY KEY,
    name            varchar(255)                NOT NULL,
    email           varchar(255)                NOT NULL,
    hashed_password char(60)                    NOT NULL,
    created         timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
