-- +goose Up
-- +goose StatementBegin
CREATE TABLE journal_generalist_embeddings (
    pmid BIGINT PRIMARY KEY REFERENCES journals(pmid),
    embeddings VECTOR(768) NOT NULL
);
CREATE INDEX ON journal_generalist_embeddings USING hnsw (embeddings vector_cosine_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS journal_generalist_embeddings;
-- +goose StatementEnd
