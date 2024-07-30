-- +goose Up
ALTER TABLE user_point_record
    ADD COLUMN scope text,
    ADD CONSTRAINT unique_user_point_record_scope UNIQUE (scope, date);

-- +goose Down
ALTER TABLE user_point_record
    DROP COLUMN scope text,
    DROP CONSTRAINT unique_user_point_record_scope;

