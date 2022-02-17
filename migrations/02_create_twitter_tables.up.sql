DROP TABLE IF EXISTS tweets CASCADE;
DROP TABLE IF EXISTS tweets_replys CASCADE;
DROP TABLE IF EXISTS tweets_likes CASCADE;
DROP TABLE IF EXISTS follows CASCADE;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE tweets
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      UUID                        NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    text         VARCHAR(280)                NOT NULL CHECK ( text <> '' ),
    image        VARCHAR(512)                CHECK ( image <> '' ),
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW()
);

CREATE TABLE tweets_replys
(
    tweet_id     BIGINT                      NOT NULL REFERENCES tweets (id) ON DELETE CASCADE,
    reply_id     BIGINT                      NOT NULL REFERENCES tweets (id) ON DELETE CASCADE,
    PRIMARY KEY(tweet_id, reply_id)
);

CREATE TABLE tweets_likes
(
    user_id      UUID                        NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    tweet_id     BIGINT                      NOT NULL REFERENCES tweets (id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, tweet_id)
);

CREATE TABLE follows
(
    follower_id  UUID                        NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    following_id UUID                        NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    PRIMARY KEY(follower_id, following_id)
);
