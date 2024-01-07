BEGIN;

    CREATE TABLE IF NOT EXISTS activityhub.status (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,

        content VARCHAR(800) NOT NULL CHECK (content <> ''), -- char limit is higher than mastodons 500 chars limit to be able to include links and html tags
        account_id UUID NOT NULL REFERENCES activityhub.account(id)
    );

    CREATE INDEX IF NOT EXISTS status_deleted_at_idx ON activityhub.status (deleted_at);

COMMIT;