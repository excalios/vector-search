package service

import (
	"context"
	"go-app/domain"
	"go-app/internal/logging"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"go.opentelemetry.io/otel"
)

type JournalRepository interface {
	GetJournalList(
		ctx context.Context,
		filter *domain.JournalFilter,
		embedding *pgvector.Vector,
	) ([]domain.JournalResponse, error)
	GetJournal(ctx context.Context, id uuid.UUID) (*domain.Journal, error)
}

type EmbeddingHTTPRepository interface {
	GetGeneralEmbedding(
		ctx context.Context,
		sentence string,
	) (*pgvector.Vector, error)
}

type JournalService struct {
	r JournalRepository
	h EmbeddingHTTPRepository
}

func NewJournalService(u JournalRepository, h EmbeddingHTTPRepository) *JournalService {
	return &JournalService{
		r: u,
		h: h,
	}
}

// GetJournal fetches a user by ID.
func (s *JournalService) GetJournal(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Journal, error) {
	tracer := otel.Tracer("service.journal")
	ctxTrace, span := tracer.Start(ctx, "JournalService.GetJournal")
	defer span.End()

	user, err := s.r.GetJournal(ctxTrace, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *JournalService) GetJournalList(
	ctx context.Context,
	filter *domain.JournalFilter,
) ([]domain.JournalResponse, error) {
	var embedding *pgvector.Vector
	var err error
	if filter != nil && filter.VSearch != "" {
		embedding, err = s.h.GetGeneralEmbedding(ctx, filter.VSearch)
		if err != nil {
			logging.LogError(ctx, err, "get_journal_list_service")
			return nil, err
		}
	}

	journals, err := s.r.GetJournalList(ctx, filter, embedding)
	if err != nil {
		logging.LogError(ctx, err, "get_journal_list_service")
		return nil, err
	}

	return journals, nil
}
