package rest_test

// import (
// 	"context"
// 	"fmt"
// 	"go-app/domain"
// 	"go-app/internal/repository/postgres"
// 	"go-app/internal/rest"
// 	"go-app/service"
// 	"net/http"
// 	"testing"
//
// 	"github.com/stretchr/testify/require"
// )
//
// func TestJournalCRUD_E2E(t *testing.T) {
// 	kit := NewTestKit(t)
//
// 	// Wire the routes and services
// 	journalRepo := postgres.NewJournalRepository(kit.DB, kit.Metrics)
// 	journalSvc := service.NewJournalService(journalRepo)
// 	rest.NewJournalHandler(kit.Echo.Group("/api/v1"), journalSvc)
//
// 	// Now start the test server
// 	kit.Start(t)
//
// 	// Get
// 	type GetType domain.ResponseSingleData[domain.Journal]
// 	getE, code := doRequest[GetType](
// 		t, http.MethodGet,
// 		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
// 		nil,
// 	)
// 	require.Equal(t, http.StatusOK, code)
// 	require.Equal(t, user.ID, getE.Data.ID)
//
// 	// Get after delete
// 	type ErrType domain.ResponseSingleData[domain.Empty]
// 	errE, code := doRequest[ErrType](
// 		t, http.MethodGet,
// 		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
// 		nil,
// 	)
// 	require.Equal(t, http.StatusNotFound, code)
// 	require.Equal(t, "Journal not found", errE.Message)
//
// 	// Hard delete, since delete API uses soft delete
// 	_, err = kit.DB.Exec(context.Background(), "DELETE from users where id = $1", user.ID)
// 	require.NoError(t, err)
// }
