package config

import (
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofiber/fiber/v2/log"
)

type MqttInstance struct {
	client mqtt.Client
}

type MqttConfig struct {
	InstanceName string
	Host         string
	Username     string
	Password     string
	ClientId     string
}

var mqttInstance = make(map[string]*MqttInstance)
var mqttMu sync.Mutex

func initializeMqtt(mqttConfig *MqttConfig) error {
	mqttMu.Lock()
	defer mqttMu.Unlock()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttConfig.Host)
	opts.SetUsername(mqttConfig.Username)
	opts.SetPassword(mqttConfig.Password)
	opts.SetClientID(mqttConfig.ClientId)

	client := mqtt.NewClient(opts)

	for i := range 5 {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Errorf("[%s] Failed to connect to MQTT broker (Attempt %d/5): %v", mqttConfig.InstanceName, i+1, token.Error())
			if i == 4 {
				return token.Error()
			}
			time.Sleep(5 * time.Second)
			continue
		}

		mqttInstance[mqttConfig.InstanceName] = &MqttInstance{
			client: client,
		}

		log.Infof("🔓 MQTT %s connection established", mqttConfig.InstanceName)
		break
	}

	return nil
}

func GetMqttInstance(instanceName string) (mqtt.Client, error) {
	mqttMu.Lock()
	defer mqttMu.Unlock()

	instance, ok := mqttInstance[instanceName]

	if !ok {
		return nil, fmt.Errorf("MQTT instance %s not found", instanceName)
	}

	return instance.client, nil
}

func CloseAllMqttInstances() {
	mqttMu.Lock()
	defer mqttMu.Unlock()

	for name, instance := range mqttInstance {
		if instance.client.IsConnectionOpen() {
			instance.client.Disconnect(250)
			log.Infof("🔒 MQTT %s connection closed", name)
		}
	}
}

func CreateMqttInstance() {
	if err := initializeMqtt(&MqttConfig{
		InstanceName: "MQTT_CLOUD",
		Host:         "tcp://hioto-rmq.pptik.id:1883",
		Username:     "/hioto:hioto",
		Password:     "ncHPk8BonsxqKyW",
		ClientId:     "listener-mqtt-cloud",
	}); err != nil {
		log.Errorf("Failed to initialize MQTT instance: %v", err)
	}

	if err := initializeMqtt(&MqttConfig{
		InstanceName: "MQTT_LOCAL",
		Host:         "tcp://127.0.0.1:1883",
		Username:     "/smarthome:smarthome",
		Password:     "Ssm4rt2!",
		ClientId:     "listener-mqtt-local",
	}); err != nil {
		log.Errorf("Failed to initialize MQTT instance: %v", err)
	}
}
