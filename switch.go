package homeassistant

import (
	"encoding/json"
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// SwitchDiscover convert sensor to mqtt format
type SwitchDiscover struct {
	UniqueID          string `json:"unique_id"`
	Name              string `json:"name"`
	StateTopic        string `json:"stat_t"`
	CommandTopic      string `json:"command_topic"`
	AvailabilityTopic string `json:"avty_t,omitempty"`
	Icon              string `json:"icon,omitempty"`
	Device            Device `json:"device,omitempty"`
}

// Switch HA sensor
type Switch struct {
	Ident           string
	Name            string
	Device          *Device
	Icon            string
	DefaultState    bool
	currentState    bool
	lastStateUpdate time.Time
	toggleFunc      func(string)
}

// NewSwitch creates a new switch with default values
func NewSwitch(ident string) Switch {
	s := Switch{
		Ident: ident,
	}
	return s
}

// GetDevice of sensor
func (s *Switch) GetDevice() *Device {
	return s.Device
}

// SetDevice of sensor
func (s *Switch) SetDevice(device *Device) {
	s.Device = device
}

// GetName of the sensor
func (s *Switch) GetName() string {
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
func (s *Switch) PublishState(broker MQTT.Client) error {
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

// SubscribeCommand subscribe to command chnnel
func (s *Switch) SubscribeCommand(broker MQTT.Client, function func(string)) error {
	broker.Subscribe(s.GetCommandTopic(), 0, s.CommandReceived)
	s.toggleFunc = function
	return nil
}

// CommandReceived when getting a message from topic
func (s *Switch) CommandReceived(broker MQTT.Client, message MQTT.Message) {
	payload := string(message.Payload())
	if payload == "ON" {
		s.toggleFunc("ON")
		s.SetState(true)
		s.PublishState(broker)
	} else {
		s.toggleFunc("OFF")
		s.SetState(false)
		s.PublishState(broker)
	}
}

// State returns current state
func (s *Switch) State() bool {
	return s.currentState
}

// SetState sets sensor state
func (s *Switch) SetState(state bool) {
	s.currentState = state
}

// lastState is the last time the sensor was updates
func (s *Switch) lastState() time.Time {
	return s.lastStateUpdate
}

// GetIdent of the sensor
func (s *Switch) GetIdent() string {
	if s.Device == nil {
		return s.Ident
	}
	return fmt.Sprintf("%s_%s", s.Device.Ident, s.Ident)
}

// GetBaseTopic for broker
func (s *Switch) GetBaseTopic() string {
	return fmt.Sprintf("homeassistant/switch/%s", s.GetIdent())
}

// GetStateTopic returns state topic
func (s *Switch) GetStateTopic() string {
	return fmt.Sprintf("%s/state", s.GetBaseTopic())
}

// GetAvailabilityTopic returns availability topic
func (s *Switch) GetAvailabilityTopic() string {
	if s.Device == nil {
		return fmt.Sprintf("%s/availability", s.GetBaseTopic())
	}
	return s.Device.GetAvailabilityTopic()
}

// PublishDiscover publish discover payload to MQTT
func (s *Switch) PublishDiscover(broker MQTT.Client) error {
	payload, err := s.GetDiscoverPayload()
	if err != nil {
		return err
	}
	token := broker.Publish(s.GetDiscoverTopic(), 0, true, payload)
	log.Infof("Publishing switch %s discovery to %s", s.GetName(), s.GetDiscoverTopic())
	log.Debug(string(payload))
	token.Wait()
	return nil
}

// GetDiscoverTopic returns discover topic
func (s *Switch) GetDiscoverTopic() string {
	return fmt.Sprintf("%s/config", s.GetBaseTopic())
}

// GetCommandTopic returns the command topic
func (s *Switch) GetCommandTopic() string {
	return fmt.Sprintf("%s/command", s.GetBaseTopic())
}

// GetDiscoverPayload generates disover payload json
func (s *Switch) GetDiscoverPayload() ([]byte, error) {
	return json.Marshal(&struct {
		UniqueID          string `json:"unique_id"`
		Name              string `json:"name"`
		StateTopic        string `json:"stat_t"`
		AvailabilityTopic string `json:"avty_t,omitempty"`
		CommandTopic      string `json:"command_topic,omitempty"`
		Icon              string `json:"icon,omitempty"`
		Device            Device `json:"device,omitempty"`
	}{
		UniqueID:          s.GetIdent(),
		Name:              s.GetName(),
		StateTopic:        s.GetStateTopic(),
		AvailabilityTopic: s.GetAvailabilityTopic(),
		CommandTopic:      s.GetCommandTopic(),
		Icon:              s.Icon,
		Device:            *s.Device,
	})
}
