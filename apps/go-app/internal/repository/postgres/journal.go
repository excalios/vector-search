package postgres

import (
	"context"
	"fmt"
	"go-app/domain"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type JournalRepository struct {
	Conn *pgxpool.Pool
}

func NewJournalRepository(conn *pgxpool.Pool) *JournalRepository {
	return &JournalRepository{
		Conn: conn,
	}
}

func (u *JournalRepository) GetJournalList(
	ctx context.Context,
	filter *domain.JournalFilter,
	embedding *pgvector.Vector,
) ([]domain.JournalResponse, error) {
	query := `
		SELECT
            pmid,
            title,
            abstract,
            content,
            mesh_terms,
            0 as distance
		FROM journals`

	if filter != nil && filter.VSearch != "" && embedding != nil {
		embeddingTable := ""
		switch filter.Type {
		case domain.GeneralVectorType:
			embeddingTable = "journal_generalist_embeddings"
		case domain.SpecialistVectorType:
			embeddingTable = "journal_specialist_embeddings"
		}
		query = fmt.Sprintf(`
            SELECT
                j.pmid,
                title,
                abstract,
                content,
                mesh_terms,
                1 - (je.embeddings <=> @query) as distance
            FROM journals j
            INNER JOIN %s je ON j.pmid = je.pmid
        `, embeddingTable)
	}

	args := pgx.StrictNamedArgs{}
	var conditions []string
	if filter != nil && filter.Search != "" {
		conditions = append(conditions, `(title ILIKE @title OR abstract ILIKE @title OR content ILIKE @title)`)
		args["title"] = "%" + filter.Search + "%"
	}

	if len(conditions) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(conditions, " AND "))
	}

	if filter != nil && filter.VSearch != "" && embedding != nil {
		query += " ORDER BY je.embeddings <-> @query "
		args["query"] = embedding
	}

	if filter.Limit != nil && filter.Page != nil {
		query += " LIMIT @limit OFFSET @offset"
		args["limit"] = *filter.Limit
		args["offset"] = *filter.Page * *filter.Limit
	}
	rows, err := u.Conn.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	journals, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.JournalResponse])
	if err != nil {
		return nil, err
	}

	return journals, nil
}

func (u *JournalRepository) GetJournal(ctx context.Context, id uuid.UUID) (*domain.Journal, error) {
	tracer := otel.Tracer("repo.journal")
	ctx, span := tracer.Start(ctx, "JournalRepository.GetJournal")
	defer span.End()

	query := `
		SELECT
            *
		FROM journals
		WHERE id = $1`

	span.SetAttributes(attribute.String("query.statement", query))
	span.SetAttributes(attribute.String("query.parameter", id.String()))
	rows, err := u.Conn.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	journal, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[domain.Journal])
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return journal, nil
}
