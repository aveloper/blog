/**
  Trigger function to set the updated_at field to current time
  when a row is updated.
 */
CREATE FUNCTION set_updated_at() RETURNS TRIGGER
    LANGUAGE plpgsql
AS
$BODY$
BEGIN
    IF TG_OP = 'UPDATE' THEN
        NEW."updated_at" = NOW();
    END IF;
    RETURN NEW;
END;
$BODY$;

/**
  Settings table
 */
CREATE TABLE settings
(
    id         SERIAL                   NOT NULL PRIMARY KEY,
    site_name  CHARACTER VARYING(200)   NOT NULL,
    favicon    CHARACTER VARYING(200)   NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_settings
    BEFORE UPDATE
    ON settings
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

/**
Type for all possible user roles
*/
CREATE TYPE user_role AS ENUM ('owner', 'admin', 'editor', 'author', 'contributor');

/**
  Users table
 */
CREATE TABLE users
(
    id             SERIAL PRIMARY KEY       NOT NULL,
    name           VARCHAR(200)             NOT NULL,
    email          VARCHAR(200)             NOT NULL UNIQUE,
    password       TEXT                     NOT NULL,
    role           user_role                NOT NULL DEFAULT 'contributor'::user_role,
    email_verified BOOLEAN                  NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE
    ON users
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();


/**
  Type for status for posts
 */
CREATE TYPE post_status AS ENUM ('draft', 'publish', 'archive');

CREATE TABLE posts
(
    id           SERIAL PRIMARY KEY       NOT NULL,
    title        VARCHAR(200)             NOT NULL,
    slug         VARCHAR(200)             NOT NULL UNIQUE,
    summary      VARCHAR(200)             NOT NULL,
    feature_img  VARCHAR(200)                      DEFAULT NULL,
    content      TEXT                     NOT NULL,
    status       CHARACTER VARYING(10)    NOT NULL,
    likes        BIGINT                   NOT NULL DEFAULT 0,
    views        BIGINT                   NOT NULL DEFAULT 0,
    published_at TIMESTAMP WITH TIME ZONE          DEFAULT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

/**
  Trigger function to set published_at to current time
  whenever a draft post is published
 */
CREATE OR REPLACE FUNCTION set_published_at() RETURNS TRIGGER
    LANGUAGE PLPGSQL
AS
$BODY$
BEGIN
    /**
      When a post is published,
      Update the published_at column to the current time
     */
    IF OLD."status" != 'publish' AND NEW."status" = 'publish' THEN
        NEW."published_at" = NOW();
    END IF;
    RETURN NEW;
END;
$BODY$;

CREATE TRIGGER set_published_at_posts
    BEFORE UPDATE
    ON posts
    FOR EACH ROW
EXECUTE FUNCTION set_published_at();

CREATE TRIGGER set_updated_at_posts
    BEFORE UPDATE
    ON posts
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

/**
  Join table for post and users
 */
CREATE TABLE post_authors
(
    user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    post_id INT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, post_id)
);

/**
  Tags table
 */
CREATE TABLE tags
(
    id         SERIAL                   NOT NULL PRIMARY KEY,
    name       VARCHAR(200)             NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_tags
    BEFORE UPDATE
    ON tags
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

/**
  Join table for posts and tags
 */
CREATE TABLE post_tags
(
    post_id INT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    tag_id  INT NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

/**
  Topics table
 */
CREATE TABLE topics
(
    id         SERIAL                   NOT NULL PRIMARY KEY,
    name       CHARACTER VARYING(200)   NOT NULL UNIQUE,
    parent_id  INT                               DEFAULT NULL REFERENCES topics (id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_topics
    BEFORE UPDATE
    ON topics
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

/**
  Join table for post and topics
 */
CREATE TABLE post_topics
(
    topic_id INT NOT NULL REFERENCES topics (id) ON DELETE CASCADE,
    post_id  INT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    PRIMARY KEY (topic_id),
    UNIQUE (topic_id, post_id)
);

/**
  Table to store subscribers
 */
CREATE TABLE subscribers
(
    id             SERIAL                   NOT NULL PRIMARY KEY,
    email          VARCHAR(200)             NOT NULL UNIQUE,
    email_verified BOOLEAN                  NOT NULL DEFAULT FALSE,
    unsubscribed   BOOLEAN                  NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at_subscribers
    BEFORE UPDATE
    ON subscribers
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
