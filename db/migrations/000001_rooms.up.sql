START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  app.rooms (
    id        uuid  NOT NULL,
    name      text  NOT NULL,
    owner_id  uuid  NOT NULL,
    movie_id  uuid,
    timecode  time,

    CONSTRAINT pk_rooms PRIMARY KEY (id)
  );

CREATE UNIQUE INDEX ix_rooms_id ON app.rooms (id);

CREATE INDEX ix_rooms_owner_id ON app.rooms (owner_id);
CREATE INDEX ix_rooms_movie_id ON app.rooms (movie_id);

COMMIT;
