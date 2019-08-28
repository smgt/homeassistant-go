package homeassistant

import "testing"

func TestSensorState(t *testing.T) {
	t.Run("Sensor state cap", func(t *testing.T) {
		s := NewSensor("sensor01")
		for i := 0; i < s.stateRetention+10; i++ {
			s.AddState(1.0)
		}
		got := len(s.States)
		want := s.stateRetention
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func TestSensorMovingAverage(t *testing.T) {
	t.Run("Test moving average", func(t *testing.T) {
		s := Sensor{}
		s.States = []float64{1.0, 1.0, 10.0, 10.0, 10.0}
		got, _ := s.MovingAverage()
		want := 6.4
		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	})

	t.Run("Test moving average cap", func(t *testing.T) {
		s := Sensor{}
		s.States = []float64{1.0, 1.0, 10.0, 10.0, 10.0, 10.0, 10.0}
		got, _ := s.MovingAverage()
		want := 10.0
		if got != want {
			t.Errorf("got %f want %f", got, want)
		}
	})

	t.Run("Test moving average without states", func(t *testing.T) {
		s := Sensor{}
		_, err := s.MovingAverage()
		if err == nil {
			t.Errorf("Didn't get error when there isn't any states")
		}
	})

}

func TestSensorTopics(t *testing.T) {
	t.Run("State topic without device", func(t *testing.T) {
		s := Sensor{Ident: "sensor1"}
		got := s.GetStateTopic()
		want := "homeassistant/sensor/sensor1/state"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("State topic with device", func(t *testing.T) {
		d := Device{Ident: "device1"}
		s := Sensor{Ident: "sensor1", Device: &d}
		got := s.GetStateTopic()
		want := "homeassistant/sensor/device1_sensor1/state"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Availability topic without device", func(t *testing.T) {
		s := Sensor{Ident: "sensor1"}
		got := s.GetAvailabilityTopic()
		want := "homeassistant/sensor/sensor1/availability"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Availability topic with device", func(t *testing.T) {
		d := Device{Ident: "device1"}
		s := Sensor{Ident: "sensor1", Device: &d}
		got := s.GetAvailabilityTopic()
		want := "device/device1/availability"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Discover topic without device", func(t *testing.T) {
		s := Sensor{Ident: "sensor1"}
		got := s.GetDiscoverTopic()
		want := "homeassistant/sensor/sensor1/config"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Discover topic with device", func(t *testing.T) {
		d := Device{Ident: "device1"}
		s := Sensor{Ident: "sensor1", Device: &d}
		got := s.GetDiscoverTopic()
		want := "homeassistant/sensor/device1_sensor1/config"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Base topic with device", func(t *testing.T) {
		d := Device{Ident: "device1"}
		s := Sensor{Ident: "sensor1", Device: &d}
		got := s.GetBaseTopic()
		want := "homeassistant/sensor/device1_sensor1"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Base topic without device", func(t *testing.T) {
		s := Sensor{Ident: "sensor1"}
		got := s.GetBaseTopic()
		want := "homeassistant/sensor/sensor1"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

}

func TestJSON(t *testing.T) {
	t.Run("Test json output", func(t *testing.T) {
		t.Skip("Not implemnted")
	})
}

func TestName(t *testing.T) {
	device := Device{Name: "Device"}
	sensor := Sensor{Name: "Sensor"}

	t.Run("Sensor name without device", func(t *testing.T) {
		got := sensor.GetName()
		want := "Sensor"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})

	t.Run("Sensor name with device", func(t *testing.T) {
		sensor.Device = &device
		got := sensor.GetName()
		want := "Device Sensor"
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}

	})
}

func TestAnomaly(t *testing.T) {
	t.Run("Basic anomaly detection with anomaly", func(t *testing.T) {
		s := NewSensor("sensor01")
		s.AnomalyDetect = true
		s.AddState(33)
		s.AddState(34)
		s.AddState(35)
		err := s.AddState(1000)
		t.Log(s.States)
		if err == nil {
			t.Errorf("Wanted error but didn't get any")
		}
	})

	t.Run("Basic anomaly detection without anomaly", func(t *testing.T) {
		s := NewSensor("sensor01")
		s.AnomalyDetect = true
		s.AddState(33)
		s.AddState(34)
		err := s.AddState(35)
		if err != nil {
			t.Errorf("Got error but didn't want one: %s", err)
		}
	})

	t.Run("Basic anomaly detection with detection off", func(t *testing.T) {
		s := Sensor{AnomalyDetect: false}
		s.AddState(33)
		s.AddState(34)
		err := s.AddState(1000)
		if err != nil {
			t.Errorf("Got error but didn't want one: %s", err)
		}
	})

}

func TestDeviceAddSensor(t *testing.T) {
	device := Device{Ident: "device01"}
	t.Run("Sensor device is present after added to device", func(t *testing.T) {
		sensor := NewSensor("sensor01")
		device.AddSensor(&sensor)
		want := &device
		got := sensor.Device
		if got != want {
			t.Errorf("got %p want %p", got, want)
		}
	})
}
