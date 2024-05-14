START TRANSACTION;

DROP INDEX IF EXISTS "app".ix_friends_user_id_1;
DROP INDEX IF EXISTS "app".ix_friends_user_id_2;
DROP TABLE IF EXISTS "app".friends;

COMMIT;

