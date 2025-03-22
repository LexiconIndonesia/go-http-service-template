package module

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adryanev/go-http-service-template/common/db"
	"github.com/adryanev/go-http-service-template/repository"
)

func TestModuleRouter(t *testing.T) {
	// Create a new module with nil DB (for testing only)
	mod := NewModule(nil)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Get the router and serve the request
	handler := mod.Router()
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := `{"status":200,"message":"Hello with dependency injection"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestModuleWithMockDB(t *testing.T) {
	// Create a mock DB
	mockDB := &MockDB{}
	// Create a new module with the mock DB
	mod := NewModule(&db.DB{
		Pool:    mockDB.Pool,
		Queries: mockDB.Queries,
	})

	// Test the module
	if mod.DB.Pool != mockDB.Pool || mod.DB.Queries != mockDB.Queries {
		t.Errorf("expected DB to be the mock DB instance")
	}
}

// MockDB implements the db.DB structure for testing
type MockDB struct {
	*db.DB
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		DB: &db.DB{
			Pool:    nil,
			Queries: &repository.Queries{},
		},
	}
}

// Ping mocks the database ping method
func (m *MockDB) Ping(ctx context.Context) error {
	return nil
}
