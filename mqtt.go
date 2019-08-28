package homeassistant

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// Broker options for the MQTT broker
type Broker struct {
	URI         string
	ClientID    string
	Username    string
	Password    string
	WillTopic   string
	WillMessage string
	client      *MQTT.Client
	opts        *MQTT.ClientOptions
	logger      *log.Entry
}

// NewBroker return a new broker
func NewBroker(b *Broker) MQTT.Client {
	opts := MQTT.NewClientOptions()
	b.logger = log.WithFields(log.Fields{"unit": "mqtt"})
	b.logger.Debugf("Connecting to MQTT server %s with client id %s", b.URI, b.ClientID)
	opts.AddBroker(b.URI)
	opts.SetClientID(b.ClientID)
	if b.Username != "" {
		opts.SetUsername(b.Username)
		opts.SetPassword(b.Password)
	}
	if b.WillTopic != "" {
		b.logger.Debugf("Adding LWT to %s with payload %s", b.WillTopic, b.WillMessage)
		opts.SetBinaryWill(b.WillTopic, []byte(b.WillMessage), 0, true)
	}
	b.opts = opts
	client := MQTT.NewClient(opts)
	return client
}
