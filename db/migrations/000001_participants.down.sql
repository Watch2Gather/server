START TRANSACTION;

DROP INDEX IF EXISTS "app".ix_participants_room_id;
DROP INDEX IF EXISTS "app".ix_participants_users_id;
DROP TABLE IF EXISTS "app".participants;

COMMIT;
