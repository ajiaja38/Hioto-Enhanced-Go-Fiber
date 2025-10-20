package messagebroker

import (
	"context"
	"go/hioto/config"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler func([]byte)

func ConsumeRmq(ctx context.Context, queueName string, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("[%s] Consumer stopped before connection", queueName)
			return
		default:
		}

		conn, err := config.RmqConnection(os.Getenv("RMQ_URI"), "Local")
		if err != nil {
			log.Errorf("[%s] Failed to establish connection: %v. Retrying...", queueName, err)
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Errorf("[%s] Failed to open channel: %v", queueName, err)
			conn.Close()
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
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Errorf("[%s] Failed to consume: %v", queueName, err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Infof("[%s] Waiting for messages ⚡️", queueName)

		// Worker pool
		jobs := make(chan []byte, 100)
		wg := &sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for body := range jobs {
					handler(body)
				}
			}()
		}

		// Consume loop
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
		ch.Close()
		conn.Close()

		if ctx.Err() != nil {
			return
		}

		log.Warnf("[%s] Reconnecting after disconnect...", queueName)
		time.Sleep(5 * time.Second)
	}
}

func ConsumeMQTTTopic(ctx context.Context, topic string, clietId string, handlerFunc MessageHandler) {
	HOST := "tcp://hioto-rmq.pptik.id:1883"
	USERNAME := "/hioto:hioto"
	PASSWORD := "ncHPk8BonsxqKyW"
	TOPIC := topic

	opts := mqtt.NewClientOptions()
	opts.AddBroker(HOST)
	opts.SetUsername(USERNAME)
	opts.SetPassword(PASSWORD)
	opts.SetClientID(clietId)

	client := mqtt.NewClient(opts)

	for i := range 5 {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("[listener-mqtt] Failed to connect to MQTT broker (Attempt %d/5): %v", i+1, token.Error())
			if i == 4 {
				return
			}
			time.Sleep(5 * time.Second)
			continue
		}
		log.Info("[listener-mqtt] Successfully connected to MQTT broker")
		break
	}

	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		go handlerFunc(msg.Payload())
	}

	if token := client.Subscribe(TOPIC, 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Errorf("Failed to subscribe: %v", token.Error())
		client.Disconnect(250)
		return
	}

	log.Infof("Subscribed to topic: %s", TOPIC)

	<-ctx.Done()

	log.Warnf("MQTT context done, cleaning up...")
	client.Unsubscribe(TOPIC)
	client.Disconnect(250)
}

func ConsumeMQTTTopicLocal(ctx context.Context, topic string, clietId string, handlerFunc MessageHandler) {
	HOST := "tcp://127.0.0.1:1883"
	USERNAME := "/smarthome:smarthome"
	PASSWORD := "Ssm4rt2!"
	TOPIC := topic

	opts := mqtt.NewClientOptions()
	opts.AddBroker(HOST)
	opts.SetUsername(USERNAME)
	opts.SetPassword(PASSWORD)
	opts.SetClientID(clietId)

	client := mqtt.NewClient(opts)

	for i := range 5 {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("[listener-mqtt] Failed to connect to MQTT broker (Attempt %d/5): %v", i+1, token.Error())
			if i == 4 {
				return
			}
			time.Sleep(5 * time.Second)
			continue
		}
		log.Infof("[listener-mqtt] Successfully connected to MQTT Local broker")
		break
	}

	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		go handlerFunc(msg.Payload())
	}

	if token := client.Subscribe(TOPIC, 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Errorf("Failed to subscribe: %v", token.Error())
		client.Disconnect(250)
		return
	}

	log.Infof("Subscribed to topic: %s -> Local", TOPIC)

	<-ctx.Done()

	log.Warnf("MQTT context done, cleaning up...")
	client.Unsubscribe(TOPIC)
	client.Disconnect(250)
}
