package handlers

import (
	"bytes"
	"encoding/json"
	"matching_system/internal/api/dto"
	"matching_system/internal/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMatchService is a mock implementation of the MatchService interface
type MockMatchService struct {
	mock.Mock
}

func (m *MockMatchService) AddSinglePersonAndMatch(req dto.AddPersonRequest) (*models.Person, []models.Match) {
	args := m.Called(req)
	return args.Get(0).(*models.Person), args.Get(1).([]models.Match)
}

func (m *MockMatchService) RemoveSinglePerson(personID string) bool {
	args := m.Called(personID)
	return args.Bool(0)
}

func (m *MockMatchService) QuerySinglePeople(limit int) []models.Person {
	args := m.Called(limit)
	return args.Get(0).([]models.Person)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewMatchHandler(t *testing.T) {
	handler := NewMatchHandler()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.matchService)
}

func TestAddSinglePersonAndMatch_Success(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.POST("/add", handler.AddSinglePersonAndMatch)

	// Test data
	requestBody := dto.AddPersonRequest{
		Name:        "Alice",
		Height:      165,
		Gender:      "female",
		WantedDates: 3,
	}

	expectedPerson := &models.Person{
		ID:          "test-id-1",
		Name:        "Alice",
		Height:      165,
		Gender:      "female",
		WantedDates: 3,
	}

	expectedMatches := []models.Match{
		{
			Person1: models.Person{ID: "test-id-2", Name: "Bob", Height: 175, Gender: "male", WantedDates: 2},
			Person2: *expectedPerson,
		},
	}

	// Mock expectations
	mockService.On("AddSinglePersonAndMatch", requestBody).Return(expectedPerson, expectedMatches)

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.AddPersonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, *expectedPerson, response.Person)
	assert.Equal(t, expectedMatches, response.Matches)
	assert.Equal(t, "person added successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestAddSinglePersonAndMatch_InvalidJSON(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.POST("/add", handler.AddSinglePersonAndMatch)

	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/add", bytes.NewBufferString(`{"name": "Alice", "height": "invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")

	// Verify service was not called
	mockService.AssertNotCalled(t, "AddSinglePersonAndMatch")
}

func TestAddSinglePersonAndMatch_ValidationError(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.POST("/add", handler.AddSinglePersonAndMatch)

	// Test data with invalid height (too low)
	requestBody := dto.AddPersonRequest{
		Name:        "Alice",
		Height:      50, // Below minimum (100)
		Gender:      "female",
		WantedDates: 3,
	}

	// Create request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/add", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")

	// Verify service was not called
	mockService.AssertNotCalled(t, "AddSinglePersonAndMatch")
}

func TestRemoveSinglePerson_Success(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.DELETE("/remove/:id", handler.RemoveSinglePerson)

	personID := "test-id-1"

	// Mock expectations
	mockService.On("RemoveSinglePerson", personID).Return(true)

	// Create request
	req, _ := http.NewRequest("DELETE", "/remove/"+personID, nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.RemovePersonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "person removed successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestRemoveSinglePerson_NotFound(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.DELETE("/remove/:id", handler.RemoveSinglePerson)

	personID := "non-existent-id"

	// Mock expectations
	mockService.On("RemoveSinglePerson", personID).Return(false)

	// Create request
	req, _ := http.NewRequest("DELETE", "/remove/"+personID, nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response dto.RemovePersonResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "person not found", response.Message)

	mockService.AssertExpectations(t)
}

func TestRemoveSinglePerson_MissingID(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.DELETE("/remove/:id", handler.RemoveSinglePerson)

	// Create request with empty ID
	req, _ := http.NewRequest("DELETE", "/remove", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Verify service was not called
	mockService.AssertNotCalled(t, "RemoveSinglePerson")
}

func TestQuerySinglePeople_Success(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.GET("/query", handler.QuerySinglePeople)

	limit := 5
	expectedPeople := []models.Person{
		{ID: "1", Name: "Alice", Height: 165, Gender: "female", WantedDates: 3},
		{ID: "2", Name: "Bob", Height: 175, Gender: "male", WantedDates: 2},
		{ID: "3", Name: "Charlie", Height: 180, Gender: "male", WantedDates: 1},
	}

	// Mock expectations
	mockService.On("QuerySinglePeople", limit).Return(expectedPeople)

	// Create request
	req, _ := http.NewRequest("GET", "/query?limit="+strconv.Itoa(limit), nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.QueryPeopleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPeople, response.People)
	assert.Equal(t, "people queried successfully", response.Message)

	mockService.AssertExpectations(t)
}

func TestQuerySinglePeople_MissingLimit(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.GET("/query", handler.QuerySinglePeople)

	// Create request without limit parameter
	req, _ := http.NewRequest("GET", "/query", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.QueryPeopleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "limit is required", response.Message)

	// Verify service was not called
	mockService.AssertNotCalled(t, "QuerySinglePeople")
}

func TestQuerySinglePeople_InvalidLimit(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.GET("/query", handler.QuerySinglePeople)

	// Create request with invalid limit parameter
	req, _ := http.NewRequest("GET", "/query?limit=invalid", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.QueryPeopleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "limit is required", response.Message)

	// Verify service was not called
	mockService.AssertNotCalled(t, "QuerySinglePeople")
}

func TestQuerySinglePeople_ZeroLimit(t *testing.T) {
	// Setup
	router := setupTestRouter()
	mockService := new(MockMatchService)
	handler := &MatchHandler{matchService: mockService}

	router.GET("/query", handler.QuerySinglePeople)

	limit := 0
	expectedPeople := []models.Person{}

	// Mock expectations
	mockService.On("QuerySinglePeople", limit).Return(expectedPeople)

	// Create request
	req, _ := http.NewRequest("GET", "/query?limit=0", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.QueryPeopleResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPeople, response.People)
	assert.Equal(t, "people queried successfully", response.Message)

	mockService.AssertExpectations(t)
}
