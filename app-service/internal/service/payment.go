package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPaymentService(rabbitmqURL string) (*PaymentService, error) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	// Declare queues
	_, err = ch.QueueDeclare("payment_requests", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare request queue: %v", err)
	}

	_, err = ch.QueueDeclare("payment_responses", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare response queue: %v", err)
	}

	return &PaymentService{
		conn:    conn,
		channel: ch,
	}, nil
}

func (s *PaymentService) Close() error {
	if err := s.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *types.PaymentRequest) (*types.PaymentResponse, error) {
	// Generate correlation ID for this request
	correlationID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Create response queue for this request
	responseQueue, err := s.channel.QueueDeclare(
		"",    // name (empty for auto-generated)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare response queue: %v", err)
	}

	// Start consuming from the response queue
	responses, err := s.channel.Consume(
		responseQueue.Name, // queue name
		"",                 // consumer
		true,               // auto-ack
		true,               // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming: %v", err)
	}

	// Marshal request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Publish request
	err = s.channel.Publish(
		"",                 // exchange
		"payment_requests", // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			ReplyTo:       responseQueue.Name,
			Body:          body,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish request: %v", err)
	}

	// Wait for response with timeout
	for {
		select {
		case msg := <-responses:
			if msg.CorrelationId != correlationID {
				continue
			}
			var response types.PaymentResponse
			if err := json.Unmarshal(msg.Body, &response); err != nil {
				return nil, fmt.Errorf("failed to unmarshal response: %v", err)
			}
			return &response, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("request timeout")
		case <-time.After(30 * time.Second):
			return nil, fmt.Errorf("request timeout")
		}
	}
}
