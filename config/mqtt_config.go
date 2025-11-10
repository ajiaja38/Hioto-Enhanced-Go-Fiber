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

	opts := mqtt.NewClientOptions().
		AddBroker(mqttConfig.Host).
		SetUsername(mqttConfig.Username).
		SetPassword(mqttConfig.Password).
		SetClientID(mqttConfig.ClientId).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(1 * time.Second)

	opts.OnConnect = func(client mqtt.Client) {
		log.Infof("🔓 MQTT %s connected successfully", mqttConfig.InstanceName)
	}

	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Errorf("⚠️ MQTT %s connection lost: %v", mqttConfig.InstanceName, err)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[%s] Failed to connect to MQTT broker: %v", mqttConfig.InstanceName, token.Error())
	}

	mqttInstance[mqttConfig.InstanceName] = &MqttInstance{
		client: client,
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
		log.Error(err)
	}

	if err := initializeMqtt(&MqttConfig{
		InstanceName: "MQTT_LOCAL",
		Host:         "tcp://127.0.0.1:1883",
		Username:     "/smarthome:smarthome",
		Password:     "Ssm4rt2!",
		ClientId:     "listener-mqtt-local",
	}); err != nil {
		log.Error(err)
	}
}
