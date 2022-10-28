package main

import (
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	c, err := LoadConfigWithDefaults("config-test.yaml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	log.Printf("%#v\n", c)
}
