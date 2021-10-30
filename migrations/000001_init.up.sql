CREATE TABLE user
(
    id           INTEGER PRIMARY KEY,
    is_anonymous bool,
    username     text,
    email        text,

    created_at   timestamp,
    updated_at   timestamp
)

CREATE TABLE headers_values
(
    id        INTEGER PRIMARY KEY,
    header_id INTEGER,
    FOREIGN KEY (header_id) REFERENCES headers (id)
)

CREATE TABLE headers
(
    id         INTEGER PRIMARY KEY,
    type       int,
    key        text,
    request_id INTEGER,
    FOREIGN KEY (request_id) REFERENCES requests_history (id)
)

CREATE TABLE requests_history
(
    id            INTEGER PRIMARY KEY,
    request_body  blob,
    response_body blob,
    user_id       INTEGER,
    created_at    timestamp,
    FOREIGN KEY (user_id) REFERENCES user (id)
)