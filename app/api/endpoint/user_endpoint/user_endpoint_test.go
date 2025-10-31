package user_endpoint

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"gocourse/app/api/endpoint/user_endpoint/user_response"
	"gocourse/common/logger"
	"gocourse/common/middleware"
	"gocourse/common/storage"
	"gocourse/common/tests"
	"gocourse/internal/domain/user_domain"
	"gocourse/internal/infrastructure/repository/user_repository"
)

func TestUserEndpoint(t *testing.T) {
	tests.SetTestEnv(t)
	db, dbCleanup := tests.CreateTestDatabase(t)
	cache, cacheCleanup := tests.CreateTestCache(t)
	app := gin.New()
	app.Use(middleware.ErrorHandler())
	app.Use(middleware.JSON())
	Register(
		app.Group("/users"), Deps{
			DB:     db,
			Cache:  storage.New(cache),
			Logger: logger.New(),
		},
	)
	server := httptest.NewServer(app)
	t.Cleanup(dbCleanup)
	t.Cleanup(cacheCleanup)
	t.Run(
		"save successfully one user", func(t *testing.T) {
			tests.FlushCache(cache)
			require.NoError(t, tests.TruncateAllTables(t.Context(), db))
			name := gofakeit.Name()
			email := gofakeit.Email()
			bodyBytes, marshalErr := json.Marshal(
				map[string]any{
					"name":  name,
					"email": email,
				},
			)
			require.NoError(t, marshalErr)
			res, postErr := http.Post(server.URL+"/users", "application/json", bytes.NewReader(bodyBytes))
			require.NoError(t, postErr)
			defer res.Body.Close()
			require.Equal(t, http.StatusCreated, res.StatusCode)
			var payload user_response.SaveOne
			require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
			user, findOneUserErr := user_repository.FindOne(t.Context(), db, payload.Id)
			require.NoError(t, findOneUserErr)
			require.Equal(t, name, user.Name)
			require.Equal(t, email, user.Email)
		},
	)
	t.Run(
		"invalid email field", func(t *testing.T) {
			tests.FlushCache(cache)
			require.NoError(t, tests.TruncateAllTables(t.Context(), db))
			name := gofakeit.Name()
			email := gofakeit.Name()
			bodyBytes, marshalErr := json.Marshal(
				map[string]any{
					"name":  name,
					"email": email,
				},
			)
			require.NoError(t, marshalErr)
			res, postErr := http.Post(server.URL+"/users", "application/json", bytes.NewReader(bodyBytes))
			require.NoError(t, postErr)
			defer res.Body.Close()
			require.Equal(t, http.StatusBadRequest, res.StatusCode)
			if res.StatusCode == http.StatusBadRequest {
				errs, readErr := io.ReadAll(res.Body)
				require.NoError(t, readErr)
				require.True(t, len(string(errs)) > 0)
			}
		},
	)
	t.Run(
		"invalid create body", func(t *testing.T) {
			tests.FlushCache(cache)
			require.NoError(t, tests.TruncateAllTables(t.Context(), db))
			res, postErr := http.Post(server.URL+"/users", "application/json", strings.NewReader("{}"))
			require.NoError(t, postErr)
			defer res.Body.Close()
			require.Equal(t, http.StatusBadRequest, res.StatusCode)
		},
	)
	t.Run(
		"find successfully one user", func(t *testing.T) {
			tests.FlushCache(cache)
			require.NoError(t, tests.TruncateAllTables(t.Context(), db))
			name := gofakeit.Name()
			email := gofakeit.Email()
			userId, createOneUserErr := user_repository.SaveOne(
				t.Context(), db, user_domain.UserEntity{
					Name:  name,
					Email: email,
				},
			)
			require.NoError(t, createOneUserErr)
			res, getErr := http.Get(server.URL + "/users/" + userId)
			require.NoError(t, getErr)
			defer res.Body.Close()
			require.Equal(t, http.StatusOK, res.StatusCode)
			var payload user_response.FindOne
			require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
			require.Equal(t, name, payload.Name)
			require.Equal(t, email, payload.Email)
		},
	)
	t.Run(
		"not found one user", func(t *testing.T) {
			tests.FlushCache(cache)
			require.NoError(t, tests.TruncateAllTables(t.Context(), db))
			res, getErr := http.Get(server.URL + "/users/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
			require.NoError(t, getErr)
			defer res.Body.Close()
			require.Equal(t, http.StatusNotFound, res.StatusCode)
		},
	)
	for _, id := range []string{"0123456789", "9876543210"} {
		t.Run(
			"invalid user:"+id, func(t *testing.T) {
				tests.FlushCache(cache)
				require.NoError(t, tests.TruncateAllTables(t.Context(), db))
				res, getErr := http.Get(server.URL + "/users/" + id)
				require.NoError(t, getErr)
				defer res.Body.Close()
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		)
	}
}
