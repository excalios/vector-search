-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE journals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pmid VARCHAR NOT NULL UNIQUE,
    title VARCHAR NOT NULL UNIQUE,
    abstract VARCHAR NOT NULL,
    content text NOT NULL,
    embeddings VECTOR(768) NOT NULL
);
CREATE INDEX ON journals USING hnsw (embeddings vector_cosine_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS journals;
-- +goose StatementEnd
