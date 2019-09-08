package homeassistant

import (
	"errors"
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// Device represents the device
type Device struct {
	Ident        string      `json:"ids,omitempty"`
	Name         string      `json:"name,omitempty"`
	Components   []Component `json:"-"`
	Manufacturer string      `json:"mf,omitempty"`
	Model        string      `json:"mdl,omitempty"`
}

// AddSensor to the device
// func (d *Device) AddSensor(sensor *Sensor) {
// 	sensor.Device = d
// 	d.Sensors = append(d.Sensors, *sensor)
// }

// AddComponent to the device
func (d *Device) AddComponent(component Component) error {
	ident := component.GetIdent()
	for _, c := range d.Components {
		if c.GetIdent() == ident {
			return fmt.Errorf("Component already added with ident %s", ident)
		}
	}
	component.SetDevice(d)
	d.Components = append(d.Components, component)
	return nil
}

// GetComponent by ident
func (d *Device) GetComponent(ident string) (Component, error) {
	for _, s := range d.Components {
		if s.GetIdent() == ident {
			return s, nil
		}
	}
	return nil, errors.New("Sensor not found")
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
