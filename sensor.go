package homeassistant

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// SensorDiscover convert sensor to mqtt format
type SensorDiscover struct {
	UniqueID          string `json:"unique_id"`
	Name              string `json:"name"`
	StateTopic        string `json:"stat_t"`
	AvailabilityTopic string `json:"avty_t,omitempty"`
	Icon              string `json:"icon,omitempty"`
	DeviceClass       string `json:"dev_cla,omitempty"`
	UnitOfMeasurement string `json:"unit_of_meas,omitempty"`
	Device            Device `json:"device,omitempty"`
}

// Sensor HA sensor
type Sensor struct {
	Ident                 string
	Name                  string
	Device                *Device
	DeviceClass           string
	Icon                  string
	UnitOfMeasurement     string
	States                []float64
	currentState          float64
	lastStateUpdate       time.Time
	AnomalyDetect         bool
	anomalyDetectFunction func(state float64) error
	stateRetention        int
}

// NewSensor creates a new sensor with default values
func NewSensor(ident string) Sensor {
	sensor := Sensor{
		stateRetention: 10,
		Ident:          ident,
	}
	return sensor
}

// GetName of the sensor
func (s *Sensor) GetName() string {
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

func (s *Sensor) basicAnomalyDetect(state float64) error {
	if len(s.States) > 0 {
		ma, err := s.MovingAverage()
		if err != nil {
			return err
		}
		// TODO: Fix dynamic configuration
		if math.Abs((state-ma)/state) > 0.1 {
			return errors.New("Change bigger than 10%")
		}
	}

	return nil
}

// AddState to the sensor
func (s *Sensor) AddState(state float64) error {
	if s.AnomalyDetect {
		err := s.basicAnomalyDetect(state)
		if err != nil {
			return err
		}
	}
	s.currentState = state
	s.lastStateUpdate = time.Now()

	newState := append(s.States, state)
	if len(newState) <= s.stateRetention {
		s.States = newState
	} else {
		s.States = newState[len(newState)-s.stateRetention:]
	}
	return nil
}

// PublishState publishes last state to broker
func (s *Sensor) PublishState(broker MQTT.Client) error {
	token := broker.Publish(s.GetStateTopic(), 0, false, fmt.Sprintf("%.1f", s.State()))
	token.Wait()
	return nil
}

// MovingAverage calculates moving average of last states
func (s *Sensor) MovingAverage() (float64, error) {
	numberOfStates := len(s.States)
	if len(s.States) == 0 {
		return 0.0, errors.New("No states to calculate MA on")
	}
	reversedStates := make([]float64, numberOfStates)
	copy(reversedStates, s.States)
	for i, j := 0, len(reversedStates)-1; i < j; i, j = i+1, j-1 {
		reversedStates[i], reversedStates[j] = reversedStates[j], reversedStates[i]
	}
	sum := 0.0
	for i, state := range reversedStates {
		if i > 4 {
			numberOfStates = 5
			break
		}
		sum += state
	}
	return sum / float64(numberOfStates), nil
}

// State returns current state
func (s *Sensor) State() float64 {
	return s.currentState
}

// lastState is the last time the sensor was updates
func (s *Sensor) lastState() time.Time {
	return s.lastStateUpdate
}

// GetIdent of the sensor
func (s *Sensor) GetIdent() string {
	if s.Device == nil {
		return s.Ident
	}
	return fmt.Sprintf("%s_%s", s.Device.Ident, s.Ident)
}

// GetBaseTopic for broker
func (s *Sensor) GetBaseTopic() string {
	return fmt.Sprintf("homeassistant/sensor/%s", s.GetIdent())
}

// GetStateTopic returns state topic
func (s *Sensor) GetStateTopic() string {
	return fmt.Sprintf("%s/state", s.GetBaseTopic())
}

// GetAvailabilityTopic returns availability topic
func (s *Sensor) GetAvailabilityTopic() string {
	if s.Device == nil {
		return fmt.Sprintf("%s/availability", s.GetBaseTopic())
	}
	return s.Device.GetAvailabilityTopic()
}

// PublishDiscover publish discover payload to MQTT
func (s *Sensor) PublishDiscover(broker MQTT.Client) error {
	payload, err := s.GetDiscoverPayload()
	if err != nil {
		return err
	}
	token := broker.Publish(s.GetDiscoverTopic(), 0, true, payload)
	log.Infof("Publishing sensor %s discovery to %s", s.GetName(), s.GetDiscoverTopic())
	log.Debug(string(payload))
	token.Wait()
	return nil
}

// GetDiscoverTopic returns discover topic
func (s *Sensor) GetDiscoverTopic() string {
	return fmt.Sprintf("%s/config", s.GetBaseTopic())
}

// GetDiscoverPayload generates disover payload json
func (s *Sensor) GetDiscoverPayload() ([]byte, error) {
	return json.Marshal(&struct {
		UniqueID          string `json:"unique_id"`
		Name              string `json:"name"`
		StateTopic        string `json:"stat_t"`
		AvailabilityTopic string `json:"avty_t,omitempty"`
		Icon              string `json:"icon,omitempty"`
		DeviceClass       string `json:"dev_cla,omitempty"`
		UnitOfMeasurement string `json:"unit_of_meas,omitempty"`
		Device            Device `json:"device,omitempty"`
	}{
		UniqueID:          s.GetIdent(),
		Name:              s.GetName(),
		StateTopic:        s.GetStateTopic(),
		AvailabilityTopic: s.GetAvailabilityTopic(),
		DeviceClass:       s.DeviceClass,
		Icon:              s.Icon,
		UnitOfMeasurement: s.UnitOfMeasurement,
		Device:            *s.Device,
	})
}

// GetDevice of sensor
func (s *Sensor) GetDevice() *Device {
	return s.Device
}

// SetDevice of sensor
func (s *Sensor) SetDevice(device *Device) {
	s.Device = device
}
