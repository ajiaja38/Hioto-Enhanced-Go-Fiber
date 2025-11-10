package router

import (
	"context"
	"go/hioto/pkg/handler/consumer"
	messagebroker "go/hioto/pkg/handler/message_broker"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type ConsumerMqtt struct {
	InstanceName, Topic string
	handlerFunc         func([]byte)
}

type ConsumerMessageBroker struct {
	consumerHandler *consumer.ConsumerHandler
	ctx             context.Context
}

func NewConsumerMessageBroker(consumerHandler *consumer.ConsumerHandler, ctx context.Context) *ConsumerMessageBroker {
	return &ConsumerMessageBroker{
		consumerHandler: consumerHandler,
		ctx:             ctx,
	}
}

func (c *ConsumerMessageBroker) StartConsumer() {
	go messagebroker.ConsumeRmq(c.ctx, os.Getenv("RMQ_HIOTO_LOCAL_INSTANCE"), os.Getenv("REGISTRATION_QUEUE"), c.consumerHandler.RegistrationFromCloudHandler)
	go messagebroker.ConsumeRmq(c.ctx, os.Getenv("RMQ_HIOTO_LOCAL_INSTANCE"), os.Getenv("RULES_QUEUE"), c.consumerHandler.RulesHandler)
	go messagebroker.ConsumeRmq(c.ctx, os.Getenv("RMQ_HIOTO_LOCAL_INSTANCE"), os.Getenv("MONITORING_QUEUE"), c.consumerHandler.MonitoringDataDevice)

	go func() {
		for {
			ctx, cancel := context.WithCancel(context.Background())
			cancels := c.startRoutingConsumer(ctx)

			time.Sleep(1 * time.Hour)

			log.Warn("⏰ Restarting routing consumers to ensure binding is refreshed")
			for _, c := range cancels {
				c()
			}
			cancel()
		}
	}()
}

func (c *ConsumerMessageBroker) startRoutingConsumer(ctx context.Context) []context.CancelFunc {
	var cancels []context.CancelFunc

	createCtx := func() context.Context {
		c, cancel := context.WithCancel(ctx)
		cancels = append(cancels, cancel)
		return c
	}

	for _, route := range []ConsumerMqtt{
		{"MQTT_CLOUD", os.Getenv("CONTROL_ROUTING_KEY"), c.consumerHandler.ControlHandler},
		{"MQTT_CLOUD", os.Getenv("REGISTRATION_ROUTING_KEY"), c.consumerHandler.RegistrationFromCloudHandler},
		{"MQTT_CLOUD", os.Getenv("UPDATE_DEVICE_ROUTING_KEY"), c.consumerHandler.UpdateDeviceFromCloudHandler},
		{"MQTT_CLOUD", os.Getenv("DELETE_DEVICE_ROUTING_KEY"), c.consumerHandler.DeleteDeviceFromCloudHandler},
		{"MQTT_LOCAL", os.Getenv("AKTUATOR_ROUTING_KEY"), c.consumerHandler.TestingConsumeAktuator},
		{"MQTT_LOCAL", os.Getenv("SENSOR_QUEUE"), c.consumerHandler.ControlSensorHandler},
	} {
		go messagebroker.ConsumeMQTTTopic(createCtx(), route.InstanceName, route.Topic, route.handlerFunc)
	}

	return cancels
}
