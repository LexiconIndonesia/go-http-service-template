package module

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adryanev/go-http-service-template/common/db"
	"github.com/adryanev/go-http-service-template/common/messaging"
	"github.com/adryanev/go-http-service-template/repository"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func TestModuleRouter(t *testing.T) {
	// Create a new module with nil dependencies (for testing only)
	mod := NewModule(nil, nil)

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
	// Create mock dependencies
	mockDB := NewMockDB()
	mockNats := NewMockNatsClient()

	// Create a new module with the mock dependencies
	mod := NewModule(mockDB, mockNats)

	// Test the module
	if mod.DB != mockDB {
		t.Errorf("expected DB to be the mock DB instance")
	}

	if mod.NatsClient != mockNats {
		t.Errorf("expected NatsClient to be the mock NATS client instance")
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

// MockNatsClient implements the messaging.NatsClient for testing
type MockNatsClient struct {
	*messaging.NatsClient
	PublishCalled bool
	SubjectsSeen  []string
}

// NewMockNatsClient creates a new mock NATS client
func NewMockNatsClient() *MockNatsClient {
	return &MockNatsClient{
		PublishCalled: false,
		SubjectsSeen:  []string{},
	}
}

// Publish mocks the publish method
func (m *MockNatsClient) Publish(subject string, data []byte) error {
	m.PublishCalled = true
	m.SubjectsSeen = append(m.SubjectsSeen, subject)
	return nil
}

// PublishAsync mocks the async publish method
func (m *MockNatsClient) PublishAsync(subject string, data []byte) (nats.PubAckFuture, error) {
	m.PublishCalled = true
	m.SubjectsSeen = append(m.SubjectsSeen, subject)
	return &MockPubAckFuture{}, nil
}

// GetStream mocks the GetStream method
func (m *MockNatsClient) GetStream(ctx context.Context, name string) (jetstream.Stream, error) {
	return &MockStream{}, nil
}

// MockPubAckFuture implements nats.PubAckFuture for testing
type MockPubAckFuture struct{}

// Ok returns a channel that's immediately ready
func (m *MockPubAckFuture) Ok() <-chan struct{} {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return ch
}

// Err returns an error channel
func (m *MockPubAckFuture) Err() <-chan error {
	return make(chan error)
}

// Sequence returns a dummy sequence number
func (m *MockPubAckFuture) Sequence() uint64 {
	return 1
}

// Stream returns a dummy stream name
func (m *MockPubAckFuture) Stream() string {
	return "MOCK_STREAM"
}

// Domain returns a dummy domain
func (m *MockPubAckFuture) Domain() string {
	return ""
}

// MockStream implements jetstream.Stream for testing
type MockStream struct{}

// Info returns mock stream info
func (m *MockStream) Info(ctx context.Context) (*jetstream.StreamInfo, error) {
	return &jetstream.StreamInfo{
		Config: jetstream.StreamConfig{
			Name: "MOCK_STREAM",
		},
	}, nil
}

// CreateOrUpdateConsumer creates a mock consumer
func (m *MockStream) CreateOrUpdateConsumer(ctx context.Context, cfg jetstream.ConsumerConfig) (jetstream.Consumer, error) {
	return &MockConsumer{}, nil
}

// Consumer returns a mock consumer
func (m *MockStream) Consumer(ctx context.Context, name string) (jetstream.Consumer, error) {
	return &MockConsumer{}, nil
}

// OrderedConsumer returns a mock ordered consumer
func (m *MockStream) OrderedConsumer(ctx context.Context, cfg jetstream.OrderedConsumerConfig) (jetstream.Consumer, error) {
	return &MockConsumer{}, nil
}

// Consumers returns a list of consumers
func (m *MockStream) Consumers(ctx context.Context) ([]jetstream.Consumer, error) {
	return []jetstream.Consumer{&MockConsumer{}}, nil
}

// ConsumerNames returns a list of consumer names
func (m *MockStream) ConsumerNames(ctx context.Context) ([]string, error) {
	return []string{"MOCK_CONSUMER"}, nil
}

// DeleteConsumer deletes a mock consumer
func (m *MockStream) DeleteConsumer(ctx context.Context, name string) error {
	return nil
}

// Purge purges the mock stream
func (m *MockStream) Purge(ctx context.Context) error {
	return nil
}

// MockConsumer implements jetstream.Consumer for testing
type MockConsumer struct{}

// Info returns mock consumer info
func (m *MockConsumer) Info(ctx context.Context) (*jetstream.ConsumerInfo, error) {
	return &jetstream.ConsumerInfo{
		Name: "MOCK_CONSUMER",
	}, nil
}

// Consume returns a mock consume context
func (m *MockConsumer) Consume(handler jetstream.MessageHandler) (jetstream.ConsumeContext, error) {
	return &MockConsumeContext{}, nil
}

// FetchNoWait fetches messages without waiting
func (m *MockConsumer) FetchNoWait(batch int) (jetstream.MessageBatch, error) {
	return &MockMessageBatch{}, nil
}

// Fetch fetches messages
func (m *MockConsumer) Fetch(batch int, timeout time.Duration) (jetstream.MessageBatch, error) {
	return &MockMessageBatch{}, nil
}

// Next gets the next message
func (m *MockConsumer) Next(timeout time.Duration) (jetstream.Msg, error) {
	return nil, nil
}

// MockConsumeContext implements jetstream.ConsumeContext for testing
type MockConsumeContext struct{}

// Stop stops the consumption
func (m *MockConsumeContext) Stop() {}

// MockMessageBatch implements jetstream.MessageBatch for testing
type MockMessageBatch struct{}

// Messages returns mock messages
func (m *MockMessageBatch) Messages() <-chan jetstream.Msg {
	return make(chan jetstream.Msg)
}

// Error returns any error
func (m *MockMessageBatch) Error() error {
	return nil
}

// Stop stops the batch
func (m *MockMessageBatch) Stop() {}
