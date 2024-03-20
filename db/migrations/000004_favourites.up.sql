START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  app.favourite (
    movie_id  uuid  NOT NULL  DEFAULT (public.uuid_generate_v4()),
    user_id   uuid  NOT NULL  DEFAULT (public.uuid_generate_v4()),

    CONSTRAINT um_users_id_fk
      FOREIGN KEY(user_id)
        REFERENCES app.users(id),
    CONSTRAINT um_movies_id_fk
      FOREIGN KEY(movie_id)
        REFERENCES app.movies(id)
  );

CREATE INDEX ix_favourite_movie_id ON app.favourite (movie_id);
CREATE INDEX ix_favourite_user_id ON app.favourite (user_id);

COMMIT;
