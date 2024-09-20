package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GuohaoMa/tucowDemo/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return db, mock
}

func TestFindPathHandler_NoQueries(t *testing.T) {
	db, _ := setupMockDB(t)
	defer db.Close()

	router := gin.Default()
	graph := &model.Graph{Db: db, Id: 2}
	router.POST("/find-path", FindPathHandler(graph))

	requestPayload := FindPathRq{}
	jsonPayload, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/find-path", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFindPathHandler_CycleDetected(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM edges WHERE graph_id = ?").
		WithArgs(3).
		WillReturnRows(sqlmock.NewRows([]string{"from_identity", "to_identity", "cost"}).
			AddRow("A", "B", 1).
			AddRow("B", "C", 1).
			AddRow("C", "A", 1))
	router := gin.Default()
	graph := &model.Graph{Db: db, Id: 3}
	router.POST("/find-path", FindPathHandler(graph))

	requestPayload := FindPathRq{
		Queries: []Query{
			{Paths: PathRq{Start: "A", End: "B"}},
		},
	}
	jsonPayload, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/find-path", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFindPathHandler_NoPathFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM edges WHERE graph_id = ?").
		WithArgs(4).
		WillReturnRows(sqlmock.NewRows([]string{"from_identity", "to_identity", "cost"}).
			AddRow("A", "B", 1))

	router := gin.Default()
	graph := &model.Graph{Db: db, Id: 4}
	router.POST("/find-path", FindPathHandler(graph))

	requestPayload := FindPathRq{
		Queries: []Query{
			{Paths: PathRq{Start: "A", End: "C"}},
		},
	}
	jsonPayload, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/find-path", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response FindPathRs
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Answers, 1)
	assert.NotNil(t, response.Answers[0].Paths)
	assert.Empty(t, response.Answers[0].Paths.AllPaths)
}

func TestFindPathHandler_NoCheapestPathFound(t *testing.T) {
	db, mock := setupMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT (.+) FROM edges WHERE graph_id = ?").
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"from_identity", "to_identity", "cost"}).
			AddRow("A", "B", 1))

	router := gin.Default()
	graph := &model.Graph{Db: db, Id: 5}
	router.POST("/find-path", FindPathHandler(graph))

	requestPayload := FindPathRq{
		Queries: []Query{
			{Cheapest: CheapestPathRq{Start: "A", End: "C"}},
		},
	}
	jsonPayload, err := json.Marshal(requestPayload)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/find-path", bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response FindPathRs
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Answers, 1)
	assert.NotNil(t, response.Answers[0].Cheapest)
	if path, ok := response.Answers[0].Cheapest.Path.(bool); ok {
		assert.False(t, path)
	} else {
		t.Errorf("Expected boolean type for Cheapest.Path, got %T", response.Answers[0].Cheapest.Path)
	}
}
