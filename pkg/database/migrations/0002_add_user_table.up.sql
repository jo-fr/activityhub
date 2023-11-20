BEGIN;
    CREATE TABLE IF NOT EXISTS activityhub.account (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,
        preferred_username VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        summary VARCHAR(255) NOT NULL,
        private_key BYTEA NOT NULL,
        public_key BYTEA NOT NULL,

        CONSTRAINT unique_preferred_username UNIQUE (preferred_username,deleted_at),
        CONSTRAINT unique_public_key UNIQUE (public_key,deleted_at),
        CONSTRAINT unique_private_key UNIQUE (private_key,deleted_at)
    );

    CREATE INDEX IF NOT EXISTS account_deleted_at_idx ON activityhub.account (deleted_at);

COMMIT;