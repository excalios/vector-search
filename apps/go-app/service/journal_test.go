package service_test

import (
	"context"
	"errors"
	"go-app/domain"
	"go-app/service"
	"go-app/service/mocks"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJournalService_GetJournal(t *testing.T) {
	mockJournalRepo := new(mocks.JournalRepository)
	journalService := service.NewJournalService(mockJournalRepo)

	ctx := context.Background()
	journalID := uuid.New()
	expectedJournal := &domain.Journal{
		ID:    journalID.String(),
		Name:  "Fetched Journal",
		Email: "fetched@example.com",
	}

	t.Run("Successfully fetches a journal", func(t *testing.T) {
		mockJournalRepo.On("GetJournal", mock.Anything, journalID).Return(expectedJournal, nil).Once()

		j, err := journalService.GetJournal(ctx, journalID)

		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Equal(t, expectedJournal.ID, j.ID)
		assert.Equal(t, expectedJournal.Name, j.Name)

		mockJournalRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockJournalRepo = new(mocks.JournalRepository)
		journalService = service.NewJournalService(mockJournalRepo)

		repoErr := errors.New("network error")
		mockJournalRepo.On("GetJournal", mock.Anything, journalID).Return(nil, repoErr).Once()

		j, err := journalService.GetJournal(ctx, journalID)

		assert.Error(t, err)
		assert.Nil(t, j)
		assert.Equal(t, repoErr, err)

		mockJournalRepo.AssertExpectations(t)
	})

	t.Run("Returns nil when journal not found in repository", func(t *testing.T) {
		mockJournalRepo = new(mocks.JournalRepository)
		journalService = service.NewJournalService(mockJournalRepo)

		mockJournalRepo.On("GetJournal", mock.Anything, journalID).Return(nil, nil).Once()

		j, err := journalService.GetJournal(ctx, journalID)

		assert.NoError(t, err)
		assert.Nil(t, j)

		mockJournalRepo.AssertExpectations(t)
	})
}

func TestJournalService_GetJournalList(t *testing.T) {
	mockJournalRepo := new(mocks.JournalRepository)
	journalService := service.NewJournalService(mockJournalRepo)

	ctx := context.Background()
	filter := &domain.JournalFilter{
		Search: "test",
	}
	expectedJournals := []domain.Journal{
		{ID: uuid.New().String(), Name: "Test Journal One"},
		{ID: uuid.New().String(), Name: "Another Test Journal"},
	}

	t.Run("Successfully fetches journals list", func(t *testing.T) {
		mockJournalRepo.On("GetJournalList", mock.Anything, filter).Return(expectedJournals, nil).Once()

		journals, err := journalService.GetJournalList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, journals)
		assert.Len(t, journals, 2)
		assert.Equal(t, expectedJournals[0].Name, journals[0].Name)

		mockJournalRepo.AssertExpectations(t)
	})

	t.Run("Returns empty list when no journals found", func(t *testing.T) {
		mockJournalRepo = new(mocks.JournalRepository)
		journalService = service.NewJournalService(mockJournalRepo)

		mockJournalRepo.On("GetJournalList", mock.Anything, filter).Return([]domain.Journal{}, nil).Once()

		journals, err := journalService.GetJournalList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, journals)
		assert.Len(t, journals, 0)

		mockJournalRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockJournalRepo = new(mocks.JournalRepository)
		journalService = service.NewJournalService(mockJournalRepo)

		repoErr := errors.New("get journals list database error")
		mockJournalRepo.On("GetJournalList", mock.Anything, filter).Return(nil, repoErr).Once()

		journals, err := journalService.GetJournalList(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, journals)
		assert.Equal(t, repoErr, err)

		mockJournalRepo.AssertExpectations(t)
	})
}
