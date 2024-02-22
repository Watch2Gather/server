START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  app.movies (
    id            uuid     NOT NULL  DEFAULT (uuid_generate_v4()),
    title         text     NOT NULL,
    description   text     NOT NULL,
    kp_rating     numeric,
    imdb_rating   numeric,
    kp_id         integer,
    year          integer  NOT NULL,
    preview       text,
    file_name     text     NOT NULL,
    country       text     NOT NULL,
    review_count  integer,

    CONSTRAINT pk_movies PRIMARY KEY (id)
  );

CREATE UNIQUE INDEX ix_movies_id ON app.movies (id);

COMMIT;
