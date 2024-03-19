START TRANSACTION;

DROP INDEX IF EXISTS "app".ix_users_id;
DROP INDEX IF EXISTS "app".ix_users_username;
DROP INDEX IF EXISTS "app".ix_users_token;
DROP TABLE IF EXISTS "app".users;

COMMIT;
