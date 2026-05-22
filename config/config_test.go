package config

import (
	"log"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	c, err := LoadConfigWithDefaults("config-test.yaml", "config-dynamic-test.yaml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	log.Printf("%#v\n", c)
}

func Test_SDRDefaults(t *testing.T) {
	c := &Config{}
	data := []byte(`debug: false
frequency: "152.00750M"
ppm: 0
router:
  url: "tls://localhost:4222"
  topic: pages
`)
	err := yaml.Unmarshal(data, c)
	if err != nil {
		t.Fatal(err)
	}
	if c.SDR.DeviceIndex != 0 {
		t.Errorf("expected default DeviceIndex 0, got %d", c.SDR.DeviceIndex)
	}
	if c.SDR.Gain != 0 {
		t.Errorf("expected default Gain 0, got %d", c.SDR.Gain)
	}
	if c.SDR.SampleRate != 22050 {
		t.Errorf("expected default SampleRate 22050, got %d", c.SDR.SampleRate)
	}
}
