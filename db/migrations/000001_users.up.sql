START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS
"app".users (
    id        uuid  NOT NULL  DEFAULT (public.uuid_generate_v4()),
    username  text  NOT NULL,
    email     text  NOT NULL,
    pwd_hash  text  NOT NULL,
    token     text,
    avatar    text,

    CONSTRAINT pk_users PRIMARY KEY (id)
  );

CREATE UNIQUE INDEX ix_users_id ON app.users (id);
CREATE UNIQUE INDEX ix_users_username ON app.users (username);
CREATE UNIQUE INDEX ix_users_token ON app.users (token);

COMMIT;
