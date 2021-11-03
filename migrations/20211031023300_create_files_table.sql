-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS files
(
    UID            CHAR(27)         NOT NULL,
    id             INT              NOT NULL AUTO_INCREMENT,
    mount_id       INT              NOT NULL, 
    user_uid      CHAR(27)          NOT NULL,
    status         INT              NULL,
    date           TIMESTAMP        NOT NULL,
    art            INT              NOT NULL,
    type           VARCHAR(255)     null,
    width          INT              NULL,
    height          INT              NULL,
    mime_type      VARCHAR(255)     NOT NULL,
    path           VARCHAR(1000)   NOT NULL,
  	filename       VARCHAR(1000)    NOT NULL,
  	org_filename   VARCHAR(1000)    NOT NULL,
  	extension      VARCHAR(100)     NOT NULL,
  	hash           VARCHAR(1000)    NOT NULL,    
    size           INT              NOT NULL,
    created_at     TIMESTAMP        NOT NULL,
    updated_at     TIMESTAMP        NULL,
    deleted_at     TIMESTAMP        NULL,
    meta           Text             NULL,
    PRIMARY KEY (UID),
    FOREIGN KEY (mount_id) REFERENCES mounts(id),
    FOREIGN KEY (user_uid) REFERENCES users(uid),
    UNIQUE KEY(id) 
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS files;