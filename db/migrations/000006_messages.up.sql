START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  app.messages (
    id          uuid                      NOT NULL  DEFAULT (uuid_generate_v4()),
    room_id     uuid                      NOT NULL  DEFAULT (uuid_generate_v4()),
    user_id     uuid                      NOT NULL,
    content     text                      NOT NULL,
    created_at  timestamp with time zone  NOT NULL  DEFAULT (now()),

    CONSTRAINT pk_messages PRIMARY KEY (id),
    CONSTRAINT msg_users_id_fk
      FOREIGN KEY(user_id)
        REFERENCES app.users(id),
    CONSTRAINT msg_rooms_id_fk
      FOREIGN KEY(room_id)
        REFERENCES app.rooms(id)
  );

CREATE UNIQUE INDEX ix_messages_id ON app.messages (id);

CREATE INDEX ix_messages_room_id ON app.messages (room_id);
CREATE INDEX ix_messages_user_id ON app.messages (user_id);

COMMIT;
