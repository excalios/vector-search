-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE journals (
    pmid BIGINT PRIMARY KEY,
    title VARCHAR NOT NULL UNIQUE,
    abstract VARCHAR NOT NULL,
    content text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS journals;
-- +goose StatementEnd
