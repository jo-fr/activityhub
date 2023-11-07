BEGIN;
    CREATE TABLE IF NOT EXISTS activityhub.follower (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP,

        account_id_followed UUID REFERENCES activityhub.account(id) NOT NULL,
        account_uri_following VARCHAR(255) NOT NULL

    );

    CREATE INDEX IF NOT EXISTS follower_deleted_at_idx ON activityhub.follower (deleted_at);

COMMIT;