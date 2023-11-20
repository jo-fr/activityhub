BEGIN;
    CREATE TABLE IF NOT EXISTS activityhub.follower (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,

        account_id_followed UUID REFERENCES activityhub.account(id) NOT NULL,
        account_uri_following VARCHAR(255) NOT NULL,

        CONSTRAINT unique_follower_account_id_followed_account_uri_following UNIQUE (deleted_at, account_id_followed, account_uri_following)

    );

    CREATE INDEX IF NOT EXISTS follower_deleted_at_idx ON activityhub.follower (deleted_at);
COMMIT;