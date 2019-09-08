package homeassistant

import MQTT "github.com/eclipse/paho.mqtt.golang"

// Component interface for all sensors etc
type Component interface {
	GetName() string
	GetIdent() string
	GetBaseTopic() string
	GetStateTopic() string
	PublishState(MQTT.Client) error
	GetAvailabilityTopic() string
	GetDiscoverTopic() string
	PublishDiscover(MQTT.Client) error
	GetDiscoverPayload() ([]byte, error)
	GetDevice() *Device
	SetDevice(*Device)
}
