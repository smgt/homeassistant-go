package homeassistant

import (
	"errors"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// Device represents the device
type Device struct {
	Ident        string   `json:"ids,omitempty"`
	Name         string   `json:"name,omitempty"`
	Sensors      []Sensor `json:"-"`
	Manufacturer string   `json:"mf,omitempty"`
	Model        string   `json:"mdl,omitempty"`
}

// AddSensor to the device
func (d *Device) AddSensor(sensor *Sensor) {
	sensor.Device = d
	d.Sensors = append(d.Sensors, *sensor)
}

// GetSensor by ident
func (d *Device) GetSensor(ident string) (Sensor, error) {
	for _, s := range d.Sensors {
		if s.Ident == ident {
			return s, nil
		}
	}
	return Sensor{}, errors.New("Sensor not found")
}

// GetAvailabilityTopic return the device availability topic for broker
func (d *Device) GetAvailabilityTopic() string {
	return fmt.Sprintf("device/%s/availability", d.Ident)
}

// PublishAvailable send availability message to broker
func (d *Device) PublishAvailable(broker MQTT.Client) error {
	token := broker.Publish(d.GetAvailabilityTopic(), 0, true, "online")
	log.Infof("Publishing device %s availability to %s", d.Name, d.GetAvailabilityTopic())
	token.Wait()
	return nil
}

// PublishUnavailable send unavailability message to broker
func (d *Device) PublishUnavailable(broker MQTT.Client) error {
	token := broker.Publish(d.GetAvailabilityTopic(), 0, true, "offline")
	log.Infof("Publishing device %s unavailability to %s", d.Name, d.GetAvailabilityTopic())
	token.Wait()
	return nil
}
