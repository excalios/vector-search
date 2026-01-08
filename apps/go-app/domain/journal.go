package domain

import (
	"github.com/pgvector/pgvector-go"
)

type Journal struct {
	PMID      int64    `json:"pmid"`
	Title     string   `json:"title"`
	Abstract  string   `json:"abstract"`
	Content   string   `json:"content"`
	MeSHTerms []string `json:"mesh_terms"`
}

type JournalEmbedding struct {
	PMID       string          `json:"pmid"`
	Embeddings pgvector.Vector `json:"embedding"`
}

type JournalResponse struct {
	PMID      int64    `json:"pmid"`
	Title     string   `json:"title"`
	Abstract  string   `json:"abstract"`
	Content   string   `json:"content"`
	MeSHTerms []string `json:"mesh_terms"`
	Distance  float64  `json:"distance"`
}

type VectorType string

const (
	GeneralVectorType    VectorType = "generalist"
	SpecialistVectorType VectorType = "specialist"
)

type JournalFilter struct {
	Limit   *int       `json:"limit" query:"limit"`
	Page    *int       `json:"page" query:"page"`
	Search  string     `json:"search" query:"search"`
	VSearch string     `json:"v_search" query:"v_search"`
	Type    VectorType `json:"type" query:"type"`
}
