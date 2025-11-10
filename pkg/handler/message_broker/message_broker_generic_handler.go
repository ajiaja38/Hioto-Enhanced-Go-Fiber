package messagebroker

import (
	"context"
	"go/hioto/config"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler func([]byte)

func ConsumeRmq(ctx context.Context, instanceName, queueName string, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("[%s] Consumer stopped before connection", queueName)
			return
		default:
		}

		instance, err := config.GetRMQInstance(instanceName)

		if err != nil {
			log.Errorf("[%s] Failed to get RabbitMQ instance: %v", queueName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		ch := instance.Channel

		if instance == nil || instance.Channel == nil {
			log.Warnf("[%s] Channel not ready, retrying...", queueName)
			time.Sleep(5 * time.Second)
			continue
		}

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
			log.Errorf("[%s] Queue declare error: %v", queueName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Errorf("[%s] Failed to consume: %v", queueName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Infof("[%s] Waiting for messages ⚡️", queueName)

		jobs := make(chan []byte, 100)
		wg := &sync.WaitGroup{}
		for range 5 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for body := range jobs {
					handler(body)
				}
			}()
		}

	consumeLoop:
		for {
			select {
			case <-ctx.Done():
				log.Warnf("[%s] Stopping consumer...", queueName)
				break consumeLoop
			case d, ok := <-msgs:
				if !ok {
					log.Warnf("[%s] Message channel closed", queueName)
					break consumeLoop
				}
				select {
				case jobs <- d.Body:
				case <-ctx.Done():
					break consumeLoop
				}
			}
		}

		close(jobs)
		wg.Wait()

		if ctx.Err() != nil {
			return
		}

		log.Warnf("[%s] Reconnecting after disconnect...", queueName)
		time.Sleep(5 * time.Second)
	}
}

func ConsumeMQTTTopic(ctx context.Context, instanceName, topic string, handlerFunc MessageHandler) {
	client, err := config.GetMqttInstance(instanceName)

	if err != nil {
		log.Error(err)
		return
	}

	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		go handlerFunc(msg.Payload())
	}

	if token := client.Subscribe(topic, 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Errorf("Failed to subscribe: %v", token.Error())
		client.Disconnect(250)
		return
	}

	log.Infof("Subscribed to topic: %s", topic)

	<-ctx.Done()

	log.Warnf("MQTT context done, cleaning up...")
	client.Unsubscribe(topic)
	client.Disconnect(250)
}
