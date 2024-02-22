START TRANSACTION;

DROP INDEX IF EXISTS "app".ix_messages_id;
DROP INDEX IF EXISTS "app".ix_messages_user_id;
DROP INDEX IF EXISTS "app".ix_messages_room_id;
DROP TABLE IF EXISTS "app".messages;

COMMIT;
