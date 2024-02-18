START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  app.participants (
    room_id   uuid  NOT NULL,
    user_id   uuid  NOT NULL,

    CONSTRAINT ur_users_id_fk
      FOREIGN KEY(user_id)
        REFERENCES app.users(id),
    CONSTRAINT ur_rooms_id_fk
      FOREIGN KEY(room_id)
        REFERENCES app.movies(id)
  );

CREATE INDEX ix_participants_room_id ON app.participants (room_id);
CREATE INDEX ix_participants_users_id ON app.participants (user_id);

COMMIT;
