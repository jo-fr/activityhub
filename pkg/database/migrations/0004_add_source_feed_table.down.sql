BEGIN;
    ALTER TABLE activityhub.account DROP COLUMN source_feed_id;
    DROP TABLE IF EXISTS activityhub.source_feed;
    DROP TYPE IF EXISTS source_feed_type;
COMMIT;