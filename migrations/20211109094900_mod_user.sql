-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE users ADD COLUMN new_password_request TIMESTAMP NULL;
ALTER TABLE users ADD COLUMN new_password_request_count INT null;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE users DROP COLUMN new_password_request_count;
ALTER TABLE users DROP COLUMN new_password_request;
