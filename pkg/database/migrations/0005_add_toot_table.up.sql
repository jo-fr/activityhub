BEGIN;

    CREATE TABLE IF NOT EXISTS activityhub.toot (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,

        content VARCHAR(500) NOT NULL CHECK (content <> ''), -- 500 char limit for content same as mastodon has
        account_id UUID NOT NULL REFERENCES activityhub.account(id)
    );

    CREATE INDEX IF NOT EXISTS follower_toot_idx ON activityhub.toot (deleted_at);

COMMIT;