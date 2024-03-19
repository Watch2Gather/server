START TRANSACTION;

DROP INDEX IF EXISTS "app".ix_rooms_id;
DROP INDEX IF EXISTS "app".ix_rooms_owner_id;
DROP INDEX IF EXISTS "app".ix_rooms_movie_id;
DROP TABLE IF EXISTS "app".rooms;

COMMIT;
