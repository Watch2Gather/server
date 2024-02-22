START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  app.users (
    id        uuid  NOT NULL  DEFAULT (uuid_generate_v4()),
    username  text  NOT NULL,
    email     text  NOT NULL,
    pwd_hash  text  NOT NULL,
    avatar    text,

    CONSTRAINT pk_users PRIMARY KEY (id)
  );

CREATE UNIQUE INDEX ix_users_id ON app.users (id);

COMMIT;
