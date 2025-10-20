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
	Key, ID     string
	handlerFunc func([]byte)
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
	go messagebroker.ConsumeRmq(c.ctx, os.Getenv("REGISTRATION_QUEUE"), c.consumerHandler.RegistrationFromCloudHandler)
	go messagebroker.ConsumeRmq(c.ctx, os.Getenv("RULES_QUEUE"), c.consumerHandler.RulesHandler)
	go messagebroker.ConsumeRmq(c.ctx, os.Getenv("MONITORING_QUEUE"), c.consumerHandler.MonitoringDataDevice)

	// Di nonaktifin dulu, ada memory leak
	// go ConsumeRmq(os.Getenv("STATUS_DEVICE_QUEUE"), c.db, log, c.consumerHandler.ChangeStatusDevice)

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

	for _, routeLocal := range []ConsumerMqtt{
		{os.Getenv("AKTUATOR_ROUTING_KEY"), "AktuatorID", c.consumerHandler.TestingConsumeAktuator},
		{os.Getenv("SENSOR_QUEUE"), "controlSensorID", c.consumerHandler.ControlSensorHandler},
	} {
		go messagebroker.ConsumeMQTTTopicLocal(createCtx(), routeLocal.Key, routeLocal.ID+os.Getenv("MAC_ADDRESS"), routeLocal.handlerFunc)
	}

	for _, route := range []ConsumerMqtt{
		{os.Getenv("CONTROL_ROUTING_KEY"), "controlID", c.consumerHandler.ControlHandler},
		{os.Getenv("REGISTRATION_ROUTING_KEY"), "registrationID", c.consumerHandler.RegistrationFromCloudHandler},
		{os.Getenv("UPDATE_DEVICE_ROUTING_KEY"), "updateID", c.consumerHandler.UpdateDeviceFromCloudHandler},
		{os.Getenv("DELETE_DEVICE_ROUTING_KEY"), "deleteID", c.consumerHandler.DeleteDeviceFromCloudHandler},
	} {
		go messagebroker.ConsumeMQTTTopic(createCtx(), route.Key, route.ID+os.Getenv("MAC_ADDRESS"), route.handlerFunc)
	}

	return cancels
}
