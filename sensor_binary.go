package homeassistant

import (
	"encoding/json"
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// BinarySensorDiscover convert sensor to mqtt format
type BinarySensorDiscover struct {
	UniqueID          string `json:"unique_id"`
	Name              string `json:"name"`
	StateTopic        string `json:"stat_t"`
	AvailabilityTopic string `json:"avty_t,omitempty"`
	Icon              string `json:"icon,omitempty"`
	DeviceClass       string `json:"dev_cla,omitempty"`
	Device            Device `json:"device,omitempty"`
}

// BinarySensor HA sensor
type BinarySensor struct {
	Ident                 string
	Name                  string
	Device                *Device
	DeviceClass           string
	Icon                  string
	currentState          bool
	lastStateUpdate       time.Time
	AnomalyDetect         bool
	anomalyDetectFunction func(state float64) error
}

// NewBinarySensor creates a new sensor with default values
func NewBinarySensor(ident string) BinarySensor {
	sensor := BinarySensor{
		Ident: ident,
	}
	return sensor
}

// GetDevice of sensor
func (s *BinarySensor) GetDevice() *Device {
	return s.Device
}

// SetDevice of sensor
func (s *BinarySensor) SetDevice(device *Device) {
	s.Device = device
}

// GetName of the sensor
func (s *BinarySensor) GetName() string {
	var name string
	if s.Name == "" {
		name = s.Ident
	} else {
		name = s.Name
	}
	if s.Device != nil {
		return fmt.Sprintf("%s %s", s.Device.Name, name)
	}
	return name
}

// PublishState publishes last state to broker
func (s *BinarySensor) PublishState(broker MQTT.Client) error {
	var state string
	if s.currentState == true {
		state = "ON"
	} else {
		state = "OFF"
	}
	token := broker.Publish(s.GetStateTopic(), 0, false, state)
	token.Wait()
	return nil
}

// State returns current state
func (s *BinarySensor) State() bool {
	return s.currentState
}

// SetState sets sensor state
func (s *BinarySensor) SetState(state bool) {
	s.currentState = state
}

// lastState is the last time the sensor was updates
func (s *BinarySensor) lastState() time.Time {
	return s.lastStateUpdate
}

// GetIdent of the sensor
func (s *BinarySensor) GetIdent() string {
	if s.Device == nil {
		return s.Ident
	}
	return fmt.Sprintf("%s_%s", s.Device.Ident, s.Ident)
}

// GetBaseTopic for broker
func (s *BinarySensor) GetBaseTopic() string {
	return fmt.Sprintf("homeassistant/binary_sensor/%s", s.GetIdent())
}

// GetStateTopic returns state topic
func (s *BinarySensor) GetStateTopic() string {
	return fmt.Sprintf("%s/state", s.GetBaseTopic())
}

// GetAvailabilityTopic returns availability topic
func (s *BinarySensor) GetAvailabilityTopic() string {
	if s.Device == nil {
		return fmt.Sprintf("%s/availability", s.GetBaseTopic())
	}
	return s.Device.GetAvailabilityTopic()
}

// PublishDiscover publish discover payload to MQTT
func (s *BinarySensor) PublishDiscover(broker MQTT.Client) error {
	payload, err := s.GetDiscoverPayload()
	if err != nil {
		return err
	}
	token := broker.Publish(s.GetDiscoverTopic(), 0, true, payload)
	log.Infof("Publishing binary sensor %s discovery to %s", s.GetName(), s.GetDiscoverTopic())
	log.Debug(string(payload))
	token.Wait()
	return nil
}

// GetDiscoverTopic returns discover topic
func (s *BinarySensor) GetDiscoverTopic() string {
	return fmt.Sprintf("%s/config", s.GetBaseTopic())
}

// GetDiscoverPayload generates disover payload json
func (s *BinarySensor) GetDiscoverPayload() ([]byte, error) {
	return json.Marshal(&struct {
		UniqueID          string `json:"unique_id"`
		Name              string `json:"name"`
		StateTopic        string `json:"stat_t"`
		AvailabilityTopic string `json:"avty_t,omitempty"`
		Icon              string `json:"icon,omitempty"`
		Device            Device `json:"device,omitempty"`
	}{
		UniqueID:          s.GetIdent(),
		Name:              s.GetName(),
		StateTopic:        s.GetStateTopic(),
		AvailabilityTopic: s.GetAvailabilityTopic(),
		Icon:              s.Icon,
		Device:            *s.Device,
	})
}
