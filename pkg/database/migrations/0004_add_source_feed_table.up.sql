BEGIN;
CREATE TYPE source_feed_type AS ENUM ('RSS');


    CREATE TABLE IF NOT EXISTS activityhub.source_feed (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ,

        name VARCHAR(255) NOT NULL CHECK (name <> ''), 
        type source_feed_type NOT NULL,
        feed_url VARCHAR(255) NOT NULL CHECK (feed_url <> ''),
        host_url VARCHAR(255) NOT NULL CHECK (host_url <> ''),
        author VARCHAR(255),
        description VARCHAR(500), -- 500 char limit for description
        image_url VARCHAR(255),
        account_id UUID NOT NULL REFERENCES activityhub.account(id) 
    );

    CREATE INDEX IF NOT EXISTS source_feed_deleted_at_idx ON activityhub.source_feed (deleted_at);
    
    CREATE UNIQUE INDEX IF NOT EXISTS unique_source_feed_name_idx ON activityhub.source_feed(name) WHERE deleted_at IS NULL;
    CREATE UNIQUE INDEX IF NOT EXISTS unique_source_feed_feed_url_idx ON activityhub.source_feed(feed_url) WHERE deleted_at IS NULL;
    CREATE UNIQUE INDEX IF NOT EXISTS unique_source_feed_account_id_idx ON activityhub.source_feed(account_id) WHERE deleted_at IS NULL;

    
COMMIT;