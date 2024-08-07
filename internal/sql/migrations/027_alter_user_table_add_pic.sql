-- +goose Up
ALTER TABLE users
    ADD COLUMN pic_no SMALLINT default 0;

-- +goose Down
ALTER TABLE users
    DROP COLUMN pic_no;