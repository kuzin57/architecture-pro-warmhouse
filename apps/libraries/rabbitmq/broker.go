package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lerenn/asyncapi-codegen/pkg/extensions"
	"github.com/streadway/amqp"
)

type RabbitMQBrokerController struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewRabbitMQBrokerController(conf *Config) (*RabbitMQBrokerController, error) {
	conn, err := amqp.Dial(conf.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RabbitMQBrokerController{
		conn:    conn,
		channel: ch,
		url:     conf.URL,
	}, nil
}

func (r *RabbitMQBrokerController) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
			return err
		}
	}
	return nil
}

func (r *RabbitMQBrokerController) Publish(ctx context.Context, channel string, mw extensions.BrokerMessage) error {
	_, err := r.channel.QueueDeclare(
		channel,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", channel, err)
	}

	log.Println("publishing message to", channel)

	err = r.channel.Publish(
		"",
		channel,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         mw.Payload,
			Headers:      convertHeaders(mw.Headers),
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to %s: %w", channel, err)
	}

	return nil
}

func (r *RabbitMQBrokerController) Subscribe(ctx context.Context, channel string) (extensions.BrokerChannelSubscription, error) {
	_, err := r.channel.QueueDeclare(
		channel,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to declare queue %s: %w", channel, err)
	}

	messages := make(chan extensions.AcknowledgeableBrokerMessage, 100)
	cancel := make(chan any, 1)

	err = r.channel.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := r.channel.Consume(
		channel,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return extensions.BrokerChannelSubscription{}, fmt.Errorf("failed to register consumer for %s: %w", channel, err)
	}

	go func() {
		defer close(messages)

		for {
			select {
			case <-ctx.Done():
				log.Printf("Context cancelled for channel %s", channel)
				return
			case <-cancel:
				log.Printf("Subscription cancelled for channel %s", channel)
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("Message channel closed for %s", channel)
					return
				}

				log.Println("received message from", channel)

				brokerMsg := extensions.BrokerMessage{
					Headers: convertAMQPHeaders(msg.Headers),
					Payload: msg.Body,
				}

				ackMsg := extensions.NewAcknowledgeableBrokerMessage(
					brokerMsg,
					&RabbitMQAcknowledgeableMessage{delivery: msg},
				)

				select {
				case messages <- ackMsg:
				case <-ctx.Done():
					msg.Nack(false, true)
					return
				}
			}
		}
	}()

	return extensions.NewBrokerChannelSubscription(messages, cancel), nil
}

func convertHeaders(headers map[string][]byte) amqp.Table {
	table := make(amqp.Table)
	for key, value := range headers {
		table[key] = string(value)
	}
	return table
}

func convertAMQPHeaders(headers amqp.Table) map[string][]byte {
	result := make(map[string][]byte)
	for key, value := range headers {
		if str, ok := value.(string); ok {
			result[key] = []byte(str)
		} else if bytes, ok := value.([]byte); ok {
			result[key] = bytes
		} else {
			if jsonBytes, err := json.Marshal(value); err == nil {
				result[key] = jsonBytes
			}
		}
	}
	return result
}

type RabbitMQAcknowledgeableMessage struct {
	delivery amqp.Delivery
}

func (m *RabbitMQAcknowledgeableMessage) AckMessage() {
	if err := m.delivery.Ack(false); err != nil {
		log.Printf("Error acknowledging message: %v", err)
	}
}

func (m *RabbitMQAcknowledgeableMessage) NakMessage() {
	if err := m.delivery.Nack(false, true); err != nil {
		log.Printf("Error nacking message: %v", err)
	}
}
