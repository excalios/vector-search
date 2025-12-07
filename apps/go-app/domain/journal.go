package domain

import (
	"github.com/pgvector/pgvector-go"
)

type Journal struct {
	ID         string          `json:"id"`
	PMID       string          `json:"pmid"`
	Title      string          `json:"title"`
	Abstract   string          `json:"abstract"`
	Content    string          `json:"content"`
	Embeddings pgvector.Vector `json:"-"`
}

type JournalResponse struct {
	ID       string  `json:"id"`
	PMID     string  `json:"pmid"`
	Title    string  `json:"title"`
	Abstract string  `json:"abstract"`
	Content  string  `json:"content"`
	Distance float64 `json:"distance"`
}

type JournalFilter struct {
	Limit   *int   `json:"limit" query:"limit"`
	Page    *int   `json:"page" query:"page"`
	Search  string `json:"search" query:"search"`
	VSearch string `json:"v_search" query:"v_search"`
}
