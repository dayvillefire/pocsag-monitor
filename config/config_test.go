package config

import (
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := LoadConfigWithDefaults("config-test.yaml", "config-dynamic-test.yaml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	log.Printf("%#v\n", c)
}
