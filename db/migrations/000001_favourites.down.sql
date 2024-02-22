START TRANSACTION;

DROP INDEX IF EXISTS "app".ix_favourite_user_id;
DROP INDEX IF EXISTS "app".ix_favourite_movie_id;
DROP TABLE IF EXISTS "app".favorites;

COMMIT;
