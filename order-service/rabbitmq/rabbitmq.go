package rabbitmq

import (
	"encoding/json"
	"log"
	"order-service/config"

	"github.com/streadway/amqp"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbit(cfg config.Config) *Client {
	url := "amqp://" + cfg.RabbitMQUser + ":" + cfg.RabbitMQPass + "@" + cfg.RabbitMQHost + ":5672/"

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	return &Client{conn: conn, channel: ch}
}

func (c *Client) Publish(queue string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = c.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = c.channel.Publish(
		"",    // exchange
		queue, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
