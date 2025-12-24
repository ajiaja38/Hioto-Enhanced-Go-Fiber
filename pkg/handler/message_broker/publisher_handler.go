package messagebroker

import (
	"go/hioto/config"

	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishToRmq(instanceName string, message []byte, queueName string, exchange string) {
	instance, err := config.GetRMQInstance(instanceName)

	if err != nil {
		log.Errorf("Failed to get RabbitMQ instance: %v 💥", err)
		return
	}

	ch := instance.Channel

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

func PublishToMqtt(instance, topic, message string) {
	client, err := config.GetMqttInstance(instance)

	if err != nil {
		log.Errorf("Failed to get MQTT instance: %v 💥", err)
		return
	}

	token := client.Publish(topic, 0, false, message)
	token.Wait()

	if token.Error() != nil {
		log.Errorf("Failed to publish message: %v 💥", token.Error())
		return
	}

	log.Infof("Published message to topic %s ✅", topic)
}
