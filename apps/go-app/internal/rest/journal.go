package rest

import (
	"context"
	"database/sql"
	"errors"
	"go-app/domain"
	"go-app/internal/logging"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type JournalService interface {
	GetJournalList(ctx context.Context, filter *domain.JournalFilter) ([]domain.JournalResponse, error)
	GetJournal(ctx context.Context, id uuid.UUID) (*domain.Journal, error)
}

type JournalHandler struct {
	Service JournalService
}

func NewJournalHandler(e *echo.Group, svc JournalService) {
	handler := &JournalHandler{
		Service: svc,
	}

	e.GET("", handler.GetJournalList)
	e.GET("/:id", handler.GetJournal)
}

// @Summary        Get Journal List
// @Description    Get All Journals
// @Tags           Journals
// @Accept         json
// @Produce        json
// @Param          filter    query        domain.JournalFilter  true "Journal filters"
// @Success        200     {object}    domain.ResponseMultipleData[domain.JournalResponse] "Successfully retrieved journal list"
// @Failure        400     {object}    domain.ResponseMultipleData[domain.Empty]              "Bad request"
// @Failure        401     {object}    domain.ResponseMultipleData[domain.Empty]              "Unauthorized"
// @Failure        500     {object}    domain.ResponseMultipleData[domain.Empty]              "Internal server error"
// @Router         /api/v1/journals [get]
func (h *JournalHandler) GetJournalList(c echo.Context) error {
	ctx := c.Request().Context()

	filter := new(domain.JournalFilter)
	if err := c.Bind(filter); err != nil {
		logging.LogWarn(ctx, "Failed to bind journal filter", slog.String("error", err.Error()))
	}

	if filter.Page == nil {
		page := 0
		filter.Page = &page
	}
	if filter.Limit == nil {
		limit := 10
		filter.Limit = &limit
	}

	journals, err := h.Service.GetJournalList(ctx, filter)
	if err != nil {
		logging.LogError(ctx, err, "get_journal_list")
		return c.JSON(http.StatusInternalServerError, domain.ResponseMultipleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to list journals: " + err.Error(),
		})
	}
	if journals == nil {
		journals = []domain.JournalResponse{}
	}

	return c.JSON(http.StatusOK, domain.ResponseMultipleData[domain.JournalResponse]{
		Data:    journals,
		Code:    http.StatusOK,
		Message: "Successfully retrieve journal list",
	})
}

// @Summary        Get Journal Detail
// @Description    Get a Journal detail
// @Tags           Journals
// @Accept         json
// @Produce        json
// @Param          id    path        string true "Journal id"
// @Success        200     {object}    domain.ResponseSingleData[domain.Journal] "Successfully retrieved journal"
// @Failure        400     {object}    domain.ResponseSingleData[domain.Empty]              "Bad request"
// @Failure        401     {object}    domain.ResponseSingleData[domain.Empty]              "Unauthorized"
// @Failure        500     {object}    domain.ResponseSingleData[domain.Empty]              "Internal server error"
// @Router         /api/v1/journals/{id} [get]
func (h *JournalHandler) GetJournal(c echo.Context) error {
	tracer := otel.Tracer("http.handler.journal")
	ctx, span := tracer.Start(c.Request().Context(), "GetJournalHandler")
	defer span.End()

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid UUID")
		return c.JSON(http.StatusBadRequest, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusBadRequest,
			Message: "Invalid journal ID format",
		})
	}

	span.SetAttributes(attribute.String("journal.id", id.String()))
	j, err := h.Service.GetJournal(ctx, id)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "not found")
			return c.JSON(http.StatusNotFound, domain.ResponseSingleData[domain.Empty]{
				Code:    http.StatusNotFound,
				Message: "Journal not found",
			})
		}

		span.SetStatus(codes.Error, "service error")
		logging.LogError(ctx, err, "get_journal")
		return c.JSON(http.StatusInternalServerError, domain.ResponseSingleData[domain.Empty]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get journal: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, domain.ResponseSingleData[domain.Journal]{
		Data:    *j,
		Code:    http.StatusOK,
		Message: "Successfully retrieved journal",
	})
}
