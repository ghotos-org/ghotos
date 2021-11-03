-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS mounts
(
    id             INT              NOT NULL AUTO_INCREMENT,
    active         INT              NULL,
    path           VARCHAR(1000)   NOT NULL,
    type           INT              NOT NULL,
    created_at     TIMESTAMP        NOT NULL,
    updated_at     TIMESTAMP        NULL,
    deleted_at     TIMESTAMP        NULL,
    PRIMARY KEY (id)
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS mounts;