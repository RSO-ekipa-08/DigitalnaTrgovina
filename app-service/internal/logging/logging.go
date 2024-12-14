package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup configures the global logger
func Setup(environment string) {
	// Set up pretty logging for development
	if environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	// Set global log level based on environment
	switch environment {
	case "development":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "staging":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Add caller information
	log.Logger = log.With().Caller().Logger()

	// Add service name
	log.Logger = log.With().Str("service", "app-service").Logger()
}

// WithRequestID adds a request ID to the logger
func WithRequestID(requestID string) zerolog.Logger {
	return log.With().Str("request_id", requestID).Logger()
}

// WithUserID adds a user ID to the logger
func WithUserID(userID string) zerolog.Logger {
	return log.With().Str("user_id", userID).Logger()
}

// WithError adds an error to the logger
func WithError(err error) *zerolog.Event {
	return log.Error().Err(err)
}
