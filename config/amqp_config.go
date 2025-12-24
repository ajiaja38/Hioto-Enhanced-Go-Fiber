package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQInstance struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

var rmqInstances = make(map[string]*RMQInstance)
var mu sync.Mutex

func initializeRabbitMQ(url, rmqInstance string) error {
	mu.Lock()
	defer mu.Unlock()

	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for i := range 5 {
		conn, err = amqp.Dial(url)

		if err == nil {
			log.Infof("✅ Successfully connected to RabbitMQ (%s)", rmqInstance)

			ch, err = conn.Channel()
			if err != nil {
				log.Infof("❌ Failed to open channel: %v", err)
				return err
			}

			rmqInstances[rmqInstance] = &RMQInstance{Conn: conn, Channel: ch}

			log.Info("✅ RabbitMQ channel opened successfully 🚀")
			return nil
		}

		log.Infof("⚠️ Failed to connect to RabbitMQ, retrying in 5 seconds... (%d/5)", i+1)
		time.Sleep(5 * time.Second)
	}

	return err
}

func GetRMQInstance(rmqtype string) (*RMQInstance, error) {
	mu.Lock()
	defer mu.Unlock()

	instance, ok := rmqInstances[rmqtype]

	if !ok {
		return nil, fmt.Errorf("RabbitMQ instance %s not found", rmqtype)
	}

	return instance, nil
}

func CloseRabbitMQ() {
	mu.Lock()
	defer mu.Unlock()

	for name, rmqInstance := range rmqInstances {
		if rmqInstance.Channel != nil {
			rmqInstance.Channel.Close()
			log.Infof("🔒 RabbitMQ %s channel closed \n", name)
		}

		if rmqInstance.Conn != nil {
			rmqInstance.Conn.Close()
			log.Infof("🔒 RabbitMQ %s connection closed", name)
		}
	}
}

func CreateRmqInstance() {
	// Init RabbitMQ Cloud
	if err := initializeRabbitMQ(RMQ_CLOUD_URI.GetValue(), RMQ_CLOUD_INSTANCE.GetValue()); err != nil {
		log.Fatal(err)
	}
}
