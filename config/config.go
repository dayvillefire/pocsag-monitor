package config

import (
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

var (
	config *Config
)

type Config struct {
	Debug        bool   `yaml:"debug" default:"false"`
	DbFile       string `yaml:"db-file" default:"scan.db"`
	Frequency    string `yaml:"frequency" default:"152.00750M"`
	PPM          int    `yaml:"ppm" default:"0"`
	ApiPort      int    `yaml:"api-port" default:"8080"`
	InstanceName string `yaml:"instance-name" default:"DEFAULT"`
	SDR          struct {
		DeviceIndex int `yaml:"device-index" default:"0"`
		Gain        int `yaml:"gain" default:"0"`
		SampleRate  int `yaml:"sample-rate" default:"22050"`
	} `yaml:"sdr"`
	Router struct {
		URL   string `yaml:"url"`
		Topic string `yaml:"topic" default:"pages"`
	} `yaml:"router"`
}

func GetConfig() *Config {
	return config
}

// UnmarshalYAML overrides default handling
func (s *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(s)

	type plain Config
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func LoadConfigWithDefaults(configPath, dynamicConfigPath string) (*Config, error) {
	c := &Config{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal([]byte(data), c)
	if err != nil {
		return c, err
	}

	config = c

	return c, err
}
