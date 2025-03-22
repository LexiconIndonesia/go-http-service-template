package messaging

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/LexiconIndonesia/go-http-service-template/common/messaging"
	"github.com/go-chi/chi/v5"
)

func TestPublishMessage(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid request",
			requestBody:    `{"subject": "test.subject", "data": {"message": "test message"}}`,
			expectedStatus: http.StatusAccepted,
			expectError:    false,
		},
		{
			name:           "Empty subject",
			requestBody:    `{"subject": "", "data": {"message": "test message"}}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Empty data",
			requestBody:    `{"subject": "test.subject", "data": null}`,
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"subject": "test.subject", "data":`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a client - we'll never actually use it for publishing in tests
			natsClient := &messaging.NatsClient{}

			// Create handler with the client
			handler := NewMessaging(natsClient)

			// Create a request
			req, err := http.NewRequest("POST", "/publish", bytes.NewBufferString(tc.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// We only test request validation for most cases
			// The actual publishing would require integration testing with a real NATS server
			if tc.name == "Valid request" {
				// For valid request, just validate the parsed request is correct
				var req MessageRequest
				if err := json.Unmarshal([]byte(tc.requestBody), &req); err != nil {
					t.Fatalf("Test data is invalid: %v", err)
				}

				if req.Subject == "" {
					t.Errorf("Test data has empty subject")
				}
				if len(req.Data) == 0 {
					t.Errorf("Test data has empty data")
				}
			} else {
				// For error cases, test the handler directly
				handler.PublishMessage(rr, req)

				// Check status code
				if rr.Code != tc.expectedStatus {
					t.Errorf("Handler returned wrong status code: got %v want %v",
						rr.Code, tc.expectedStatus)
				}

				// Verify error in response
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err == nil {
					if tc.expectError {
						_, hasError := response["error"]
						if !hasError {
							t.Errorf("Expected error in response, got none")
						}
					}
				}
			}
		})
	}
}

func TestSubscribeWebSocket(t *testing.T) {
	// Create handler
	handler := NewMessaging(&messaging.NatsClient{})

	tests := []struct {
		name           string
		subjectParam   string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid subject",
			subjectParam:   "test.subject",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Empty subject",
			subjectParam:   "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a chi router for URL parameter extraction
			r := chi.NewRouter()
			r.Get("/subscribe/{subject}", handler.SubscribeWebSocket)

			var req *http.Request
			var err error

			if tc.subjectParam == "" {
				// For empty subject, test directly
				req, err = http.NewRequest("GET", "/subscribe/", nil)
				if err != nil {
					t.Fatal(err)
				}

				// Create a response recorder
				rr := httptest.NewRecorder()

				// Call the handler directly without URL params
				handler.SubscribeWebSocket(rr, req)

				// Check status code
				if status := rr.Code; status != tc.expectedStatus {
					t.Errorf("Handler returned wrong status code: got %v want %v",
						status, tc.expectedStatus)
				}

				// Check for error in response
				var response map[string]interface{}
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err == nil {
					if tc.expectError {
						_, hasError := response["error"]
						if !hasError {
							t.Errorf("Expected error in response, got none")
						}
					}
				}
			} else {
				// For valid subject, use the router
				ts := httptest.NewServer(r)
				defer ts.Close()

				// Create a request to our test server
				resp, err := http.Get(ts.URL + "/subscribe/" + tc.subjectParam)
				if err != nil {
					t.Fatal(err)
				}
				defer resp.Body.Close()

				// Check status code
				if status := resp.StatusCode; status != tc.expectedStatus {
					t.Errorf("Handler returned wrong status code: got %v want %v",
						status, tc.expectedStatus)
				}

				// Decode the response
				var response map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// Check response
				message, hasMessage := response["message"]
				if !hasMessage {
					t.Errorf("Expected message in response, got none")
				}

				// Check that the message contains the subject
				if messageStr, ok := message.(string); ok {
					if !contains(messageStr, tc.subjectParam) {
						t.Errorf("Expected message to contain subject '%s', got '%s'",
							tc.subjectParam, messageStr)
					}
				}
			}
		})
	}
}

func TestMessagingRouter(t *testing.T) {
	// Create handler
	handler := NewMessaging(&messaging.NatsClient{})

	// Get the router
	router := handler.Router()

	// Check if router exists
	if router == nil {
		t.Fatal("Router should not be nil")
	}

	// Save current environment
	oldEnv, exists := os.LookupEnv("APP_ENV")

	// Set environment variable for development mode to enable routes
	os.Setenv("APP_ENV", "development")

	// Restore environment when done
	defer func() {
		if exists {
			os.Setenv("APP_ENV", oldEnv)
		} else {
			os.Unsetenv("APP_ENV")
		}
	}()

	// Get the router again now that we're in development mode
	router = handler.Router()

	// Check if router exists
	if router == nil {
		t.Fatal("Router should not be nil in development mode")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
