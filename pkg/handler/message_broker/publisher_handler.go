package messagebroker

import (
	"go/hioto/config"

	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishToRmq(URI string, rmqType string, message []byte, queueName string, exchange string) {
	conn, err := config.RmqConnection(URI, rmqType)

	if err != nil {
		log.Errorf("Failed to establish connection: %v 💥", err)
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Errorf("Failed to open channel: %v 💥", err)
		return
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": int32(120000),
		},
	)

	if err != nil {
		log.Errorf("Failed to declare queue: %v 💥", err)
		return
	}

	err = ch.Publish(
		exchange,
		q.Name,
		false,
		false,
		amqp.Publishing{
			Body: message,
		},
	)

	if err != nil {
		log.Errorf("Failed to publish message: %v 💥", err)
		return
	}

	log.Infof("Published message to queue %s ✅", queueName)
}

func PublishToRoutingKey(URI string, rmqType string, message []byte, exchange string, routingKey string) {
	conn, err := config.RmqConnection(URI, rmqType)
	if err != nil {
		log.Errorf("Failed to establish connection: %v 💥", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("Failed to open channel: %v 💥", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Errorf("Failed to declare exchange: %v 💥", err)
		return
	}

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Body: message,
		},
	)

	if err != nil {
		log.Errorf("Failed to publish message: %v 💥", err)
		return
	}

	log.Infof("Published message to routingKey %s ✅", routingKey)
}
