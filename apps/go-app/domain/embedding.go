package domain

import "github.com/pgvector/pgvector-go"

type EmbeddingInput struct {
	Sentence string     `json:"sentence"`
	Type     VectorType `json:"type"`
}

type EmbeddingOutput struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    pgvector.Vector `json:"data"`
}
