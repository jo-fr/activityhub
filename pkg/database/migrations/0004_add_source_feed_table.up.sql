BEGIN;
CREATE TYPE source_feed_type AS ENUM ('RSS');


    CREATE TABLE IF NOT EXISTS activityhub.source_feed (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,

        name VARCHAR(255) NOT NULL CHECK (name <> ''), 
        type source_feed_type NOT NULL,
        url VARCHAR(255) NOT NULL CHECK (url <> ''),
        description VARCHAR(500) NOT NULL CHECK (description <> '') -- 500 char limit for description
    );

    CREATE INDEX IF NOT EXISTS source_feed_deleted_at_idx ON activityhub.source_feed (deleted_at);
    
    CREATE UNIQUE INDEX IF NOT EXISTS unique_source_feed_name_idx ON activityhub.source_feed(name) WHERE deleted_at IS NULL;
    CREATE UNIQUE INDEX IF NOT EXISTS unique_source_feed_url_idx ON activityhub.source_feed(url) WHERE deleted_at IS NULL;

    
    -- add source_feed_id to account
    ALTER TABLE activityhub.account 
        ADD COLUMN source_feed_id UUID REFERENCES activityhub.source_feed(id);
    
    CREATE UNIQUE INDEX unique_account_source_feed_id_idx ON activityhub.account(source_feed_id) WHERE deleted_at IS NULL;
COMMIT;