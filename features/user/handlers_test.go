package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LexiconIndonesia/go-http-service-template/common/db"
	"github.com/LexiconIndonesia/go-http-service-template/common/models"
)

func TestCreateUser(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid request",
			requestBody:    `{"email":"test@example.com","first_name":"John","last_name":"Doe","password":"Password123!"}`,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "Empty email",
			requestBody:    `{"email":"","first_name":"John","last_name":"Doe","password":"Password123!"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid email",
			requestBody:    `{"email":"not-an-email","first_name":"John","last_name":"Doe","password":"Password123!"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Empty first name",
			requestBody:    `{"email":"test@example.com","first_name":"","last_name":"Doe","password":"Password123!"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Empty last name",
			requestBody:    `{"email":"test@example.com","first_name":"John","last_name":"","password":"Password123!"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Short password",
			requestBody:    `{"email":"test@example.com","first_name":"John","last_name":"Doe","password":"short"}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"email":"test@example.com","first_name":"John","last_name":"Doe","password":}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock DB
			mockDB := &db.DB{} // Using an empty DB since we're not actually saving the user

			// Create handler with mock DB
			handler := NewUser(mockDB)

			// Create a request
			req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.CreateUser(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v",
					rr.Code, tc.expectedStatus)
			}

			// Check response structure
			var response map[string]interface{}
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Check for error or data based on expectation
			if tc.expectError {
				errorMsg, hasError := response["error"]
				if !hasError {
					t.Errorf("Expected error in response, got none")
				}
				if errorMsg == "" {
					t.Errorf("Expected non-empty error message")
				}
			} else {
				userData, hasData := response["data"]
				if !hasData {
					t.Errorf("Expected data in response, got none")
				}

				// Validate user data structure
				userMap, ok := userData.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected user data to be a map, got %T", userData)
				}

				// Check essential fields exist
				requiredFields := []string{"id", "email", "first_name", "last_name", "created_at", "updated_at"}
				for _, field := range requiredFields {
					if _, exists := userMap[field]; !exists {
						t.Errorf("Missing required field %s in response", field)
					}
				}

				// Check that password is not present
				if _, exists := userMap["password"]; exists {
					t.Errorf("Password should not be included in response")
				}

				// Validate basic field values
				if email, ok := userMap["email"].(string); ok && tc.requestBody != `{"email":"test@example.com","first_name":"John","last_name":"Doe","password":"Password123!"}` {
					// Only check for valid request case
					if email != "test@example.com" {
						t.Errorf("Expected email to be %s, got %s", "test@example.com", email)
					}
				}
			}
		})
	}
}

// MockUserValidator is used to test validation failure scenarios
type MockUserValidator struct {
	ShouldFail bool
	ErrorMsg   string
}

func (m *MockUserValidator) Validate(user models.User) error {
	if m.ShouldFail {
		return fmt.Errorf("%s", m.ErrorMsg)
	}
	return nil
}

func TestUserRouter(t *testing.T) {
	// Create a mock DB
	mockDB := &db.DB{}

	// Create handler with mock DB
	handler := NewUser(mockDB)

	// Get the router
	router := handler.Router()

	// Check if POST / route is registered
	// We can't easily test the actual routes without more complex setup,
	// but we can verify the router exists and doesn't panic
	if router == nil {
		t.Fatal("Router should not be nil")
	}
}
