-- +goose Up
ALTER TABLE users
ADD COLUMN first_name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN last_name VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE users
DROP COLUMN first_name,
DROP COLUMN last_name; 