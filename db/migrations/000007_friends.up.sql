START TRANSACTION;

CREATE SCHEMA IF NOT EXISTS "app";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE
  app.friends (
    user_id_1  uuid  NOT NULL,
    user_id_2  uuid  NOT NULL,

    CONSTRAINT fr_users_id_1_fk
      FOREIGN KEY(user_id_1)
        REFERENCES app.users(id),
    CONSTRAINT fr_users_id_2_fk
      FOREIGN KEY(user_id_2)
        REFERENCES app.users(id)
  );

CREATE INDEX ix_friends_user_id_1 ON app.friends (user_id_1);
CREATE INDEX ix_friends_user_id_2 ON app.friends (user_id_2);

COMMIT;

