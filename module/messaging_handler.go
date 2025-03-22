package module

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adryanev/go-http-service-template/common/messaging"
	"github.com/adryanev/go-http-service-template/common/utils"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

// MessageRequest represents a request to publish a message
type MessageRequest struct {
	Subject string          `json:"subject" validate:"required"`
	Data    json.RawMessage `json:"data" validate:"required"`
}

// publishMessage handles publishing a message to NATS
func (m *Module) publishMessage(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Subject == "" {
		utils.WriteError(w, http.StatusBadRequest, "Subject is required")
		return
	}

	if len(req.Data) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "Data is required")
		return
	}

	// Check if NATS client is available
	if m.NatsClient == nil {
		log.Error().Msg("NATS client is not available")
		utils.WriteError(w, http.StatusInternalServerError, "Messaging service is not available")
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check if we need to create a stream for this subject
	// This would normally be done during service setup, but for demo we'll do it here
	streamName := "MESSAGES"
	stream, err := ensureStream(ctx, m.NatsClient, streamName, []string{req.Subject, fmt.Sprintf("%s.*", req.Subject)})
	if err != nil {
		log.Error().Err(err).Msg("Failed to ensure stream exists")
		utils.WriteError(w, http.StatusInternalServerError, "Failed to ensure messaging infrastructure")
		return
	}

	// Publish message to JetStream
	ack, err := m.NatsClient.PublishAsync(req.Subject, req.Data)
	if err != nil {
		log.Error().Err(err).Str("subject", req.Subject).Msg("Failed to publish message")
		utils.WriteError(w, http.StatusInternalServerError, "Failed to publish message")
		return
	}

	// Wait for acknowledgement with a timeout
	select {
	case <-ack.Ok():
		// Message was successfully stored
		info := struct {
			Stream   string `json:"stream"`
			Sequence uint64 `json:"sequence"`
			Subject  string `json:"subject"`
		}{
			Stream:   ack.Stream(),
			Sequence: ack.Sequence(),
			Subject:  req.Subject,
		}
		utils.WriteJSON(w, http.StatusAccepted, info)
	case err := <-ack.Err():
		// There was an error
		log.Error().Err(err).Str("subject", req.Subject).Msg("Failed to get acknowledgement for message")
		utils.WriteError(w, http.StatusInternalServerError, "Failed to confirm message delivery")
	case <-ctx.Done():
		// Timeout
		log.Warn().Str("subject", req.Subject).Msg("Timeout waiting for message acknowledgement")
		utils.WriteError(w, http.StatusRequestTimeout, "Timeout waiting for message confirmation")
	}
}

// subscribeWebSocket sets up a WebSocket connection to receive messages from a NATS subject
func (m *Module) subscribeWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get subject from URL
	subject := chi.URLParam(r, "subject")
	if subject == "" {
		utils.WriteError(w, http.StatusBadRequest, "Subject is required")
		return
	}

	// In a real implementation, this would set up a WebSocket connection
	// and stream messages from NATS to the client
	// For this example, we'll just return a message indicating what would happen

	utils.WriteMessage(w, http.StatusOK, fmt.Sprintf(
		"In a complete implementation, this would open a WebSocket connection and stream messages from subject '%s'",
		subject))
}

// ensureStream ensures that a stream exists with the given name and subjects
func ensureStream(ctx context.Context, client *messaging.NatsClient, name string, subjects []string) (jetstream.Stream, error) {
	// Try to get the stream first
	stream, err := client.GetStream(ctx, name)
	if err == nil {
		// Stream exists
		return stream, nil
	}

	// Create the stream
	streamConfig := jetstream.StreamConfig{
		Name:     name,
		Subjects: subjects,
		Storage:  jetstream.MemoryStorage,
		MaxAge:   24 * time.Hour, // Messages expire after 24 hours
	}

	return client.CreateStream(ctx, streamConfig)
}
