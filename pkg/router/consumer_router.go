package router

import (
	"context"
	"fmt"
	"go/hioto/config"
	"go/hioto/pkg/handler/consumer"
	messagebroker "go/hioto/pkg/handler/message_broker"

	"github.com/gofiber/fiber/v2/log"
)

type ConsumerMqtt struct {
	InstanceName string
	Topic        string
	HandlerFunc  func([]byte)
}

type ConsumerMessageBroker struct {
	ctx             context.Context
	consumerHandler *consumer.ConsumerHandler
}

func NewConsumerMessageBroker(
	ctx context.Context,
	consumerHandler *consumer.ConsumerHandler,
) *ConsumerMessageBroker {
	return &ConsumerMessageBroker{
		ctx:             ctx,
		consumerHandler: consumerHandler,
	}
}

func (c *ConsumerMessageBroker) StartConsumer() {
	routes := []ConsumerMqtt{
		{
			InstanceName: config.MQTT_CLOUD_INSTANCE_NAME.GetValue(),
			Topic: fmt.Sprintf(
				"%s/%s",
				config.CONTROL_ROUTING_KEY.GetValue(),
				config.MAC_ADDRESS.GetValue(),
			),
			HandlerFunc: c.consumerHandler.ControlHandler,
		},
		{
			InstanceName: config.MQTT_CLOUD_INSTANCE_NAME.GetValue(),
			Topic: fmt.Sprintf(
				"%s/%s",
				config.REGISTRATION_ROUTING_KEY.GetValue(),
				config.MAC_ADDRESS.GetValue(),
			),
			HandlerFunc: c.consumerHandler.RegistrationFromCloudHandler,
		},
		{
			InstanceName: config.MQTT_CLOUD_INSTANCE_NAME.GetValue(),
			Topic: fmt.Sprintf(
				"%s/%s",
				config.UPDATE_DEVICE_ROUTING_KEY.GetValue(),
				config.MAC_ADDRESS.GetValue(),
			),
			HandlerFunc: c.consumerHandler.UpdateDeviceFromCloudHandler,
		},
		{
			InstanceName: config.MQTT_CLOUD_INSTANCE_NAME.GetValue(),
			Topic: fmt.Sprintf(
				"%s/%s",
				config.DELETE_DEVICE_ROUTING_KEY.GetValue(),
				config.MAC_ADDRESS.GetValue(),
			),
			HandlerFunc: c.consumerHandler.DeleteDeviceFromCloudHandler,
		},
		{
			InstanceName: config.MQTT_LOCAL_INSTANCE_NAME.GetValue(),
			Topic:        config.AKTUATOR_TOPIC.GetValue(),
			HandlerFunc:  c.consumerHandler.TestingConsumeAktuator,
		},
		{
			InstanceName: config.MQTT_LOCAL_INSTANCE_NAME.GetValue(),
			Topic:        config.SENSOR_TOPIC.GetValue(),
			HandlerFunc:  c.consumerHandler.ControlSensorHandler,
		},
		{
			InstanceName: config.MQTT_LOCAL_INSTANCE_NAME.GetValue(),
			Topic:        config.MONITORING_TOPIC.GetValue(),
			HandlerFunc:  c.consumerHandler.MonitoringDataDevice,
		},
	}

	for _, route := range routes {
		go messagebroker.ConsumeMQTTTopic(
			c.ctx,
			route.InstanceName,
			route.Topic,
			route.HandlerFunc,
		)
	}

	log.Info("✅ MQTT consumers started successfully")
}
