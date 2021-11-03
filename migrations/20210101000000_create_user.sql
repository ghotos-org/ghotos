-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS users
(
    uid            CHAR(27)     NOT NULL,
    id             INT          NOT NULL AUTO_INCREMENT,
    email          VARCHAR(255) NOT NULL,
    password       VARCHAR(255) NOT NULL,
    created_at     TIMESTAMP    NOT NULL,
    updated_at     TIMESTAMP    NULL,
    deleted_at     TIMESTAMP    NULL,
    PRIMARY KEY (UID),
    UNIQUE KEY(id) 
);

INSERT INTO users (uid,email,password,created_at) VALUES("200KPME8UZ4tjpP0IqMyyAizmsy", "system@bluffy.de","mgr", NOW());
COMMIT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;