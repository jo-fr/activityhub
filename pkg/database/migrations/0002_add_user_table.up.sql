BEGIN;
    CREATE TABLE IF NOT EXISTS activityhub.account (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMP,
        preferred_username VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        summary VARCHAR(255) NOT NULL,
        private_key BYTEA NOT NULL,
        public_key BYTEA NOT NULL
    );

    CREATE INDEX IF NOT EXISTS account_deleted_at_idx ON activityhub.account (deleted_at);
    CREATE UNIQUE INDEX IF NOT EXISTS account_preferred_username_idx ON activityhub.account (preferred_username,deleted_at);
    CREATE UNIQUE INDEX IF NOT EXISTS account_public_key_idx ON activityhub.account (public_key,deleted_at);
    CREATE UNIQUE INDEX IF NOT EXISTS account_private_key_idx ON activityhub.account (private_key,deleted_at);
COMMIT;