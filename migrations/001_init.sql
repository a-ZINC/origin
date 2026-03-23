CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id        UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    email     TEXT      UNIQUE NOT NULL,
    username  TEXT,
    password  TEXT      NOT NULL,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS series (
    id          UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT      UNIQUE NOT NULL,
    slug        TEXT      UNIQUE NOT NULL,
    description TEXT,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id            UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    title         TEXT      NOT NULL,
    slug          TEXT      UNIQUE NOT NULL,
    excerpt       TEXT      NOT NULL,
    body          TEXT      NOT NULL,
    "coverImage"  TEXT,
    tags          TEXT[]    NOT NULL DEFAULT '{}',
    published     BOOLEAN   NOT NULL DEFAULT FALSE,
    visibility    TEXT      NOT NULL DEFAULT 'public',
    "readingTime" INTEGER   NOT NULL DEFAULT 0,
    "seriesId"    UUID,
    "seriesPos"   INTEGER,
    "authorId"    UUID,
    "createdAt"   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updatedAt"   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY ("authorId") REFERENCES users(id),
    FOREIGN KEY ("seriesId") REFERENCES series(id)
);

CREATE TABLE IF NOT EXISTS post_publications (
    id              UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    "postId"        UUID      NOT NULL,
    platform        TEXT      NOT NULL,
    "platformPostId" TEXT,
    "platformUrl"   TEXT,
    "canonicalUrl"  TEXT      NOT NULL,
    "syncStatus"    TEXT      NOT NULL DEFAULT 'unsynced',
    "syncedAt"      TIMESTAMP,
    "createdAt"     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY ("postId") REFERENCES posts(id),
    UNIQUE(platform, "postId")
);

CREATE TABLE IF NOT EXISTS platform_analtics (
    id              UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    "publicationId" UUID      NOT NULL,
    views           INTEGER   NOT NULL DEFAULT 0,
    reactions       INTEGER   NOT NULL DEFAULT 0,
    comments        INTEGER   NOT NULL DEFAULT 0,
    bookmarks       INTEGER   NOT NULL DEFAULT 0,
    "fetchedAt"     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY ("publicationId") REFERENCES post_publications(id)
);

CREATE TABLE IF NOT EXISTS site_analytics (
    "postId"         UUID    NOT NULL,
    views            INTEGER NOT NULL DEFAULT 0,
    "uniqueViews"    INTEGER NOT NULL DEFAULT 0,
    "avgReadTimes"   INTEGER NOT NULL DEFAULT 0,
    "avgScrollDepth" INTEGER NOT NULL DEFAULT 0,
    date             DATE    NOT NULL DEFAULT CURRENT_DATE,

    FOREIGN KEY ("postId") REFERENCES posts(id),
    PRIMARY KEY ("postId", date)
);

CREATE TABLE IF NOT EXISTS projects (
    id          UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT      NOT NULL,
    description TEXT      NOT NULL,
    url         TEXT,
    "repoUrl"   TEXT,
    language    TEXT,
    stars       INTEGER   NOT NULL DEFAULT 0,
    featured    BOOLEAN   NOT NULL DEFAULT FALSE,
    "sortOrder" INTEGER   NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_slug      ON posts (slug);
CREATE INDEX IF NOT EXISTS idx_posts_series    ON posts ("seriesId");
CREATE INDEX IF NOT EXISTS idx_posts_published ON posts (published, visibility);
CREATE INDEX IF NOT EXISTS idx_posts_tags      ON posts USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_pub_posts       ON post_publications ("postId");
CREATE INDEX IF NOT EXISTS idx_post_analytics  ON post_publications ("postId");
CREATE INDEX IF NOT EXISTS idx_platform_analytics ON platform_analtics ("publicationId");
CREATE INDEX IF NOT EXISTS idx_site_analytics  ON site_analytics ("postId");