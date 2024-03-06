-- +goose Up
CREATE TABLE IF NOT EXISTS user_trackedmetric (
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    metric_id bigint NOT NULL REFERENCES trackedmetrics ON DELETE CASCADE,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_trackedmetric UNIQUE (user_id, metric_id),
    PRIMARY KEY (user_id, metric_id, granted_at)
);

CREATE TABLE IF NOT EXISTS user_smiley (
    user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    smiley_id bigint NOT NULL REFERENCES smiley ON DELETE CASCADE,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    tags text[] NOT NULL,
    PRIMARY KEY (user_id, smiley_id, granted_at)
);


-- +goose Down
DROP TABLE IF EXISTS user_trackedmetric;
DROP TABLE IF EXISTS user_smiley;