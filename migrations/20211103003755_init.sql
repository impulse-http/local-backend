-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE user
(
    id           INTEGER PRIMARY KEY,
    is_anonymous bool,
    username     text,
    email        text,

    created_at   timestamp,
    updated_at   timestamp
);

CREATE TABLE headers_values
(
    id           INTEGER PRIMARY KEY,
    header_id    INTEGER,
    header_value text,
    FOREIGN KEY (header_id) REFERENCES headers (id)
);

CREATE TABLE headers
(
    id         INTEGER PRIMARY KEY,
    key        text,
    request_id INTEGER,
    is_request INTEGER,
    FOREIGN KEY (request_id) REFERENCES requests_history (id)
);

CREATE TABLE requests_history
(
    id            INTEGER PRIMARY KEY,
    request_body  blob,
    response_body blob,
    user_id       INTEGER,
    created_at    timestamp,
    method        string,
    FOREIGN KEY (user_id) REFERENCES user (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE user;
DROP TABLE requests_history;
DROP TABLE headers_values;
DROP TABLE headers;
-- +goose StatementEnd
