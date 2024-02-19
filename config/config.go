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
	Debug          bool           `yaml:"debug" default:"false"`
	DbFile         string         `yaml:"db-file" default:"scan.db"`
	RtlFmBinary    string         `yaml:"rtlfm" default:"rtl_fm"`
	Frequency      string         `yaml:"frequency" default:"152.00750M"`
	MultiMonBinary string         `yaml:"multimon" default:"multimon-ng"`
	PPM            int            `yaml:"ppm" default:"0"`
	ApiPort        int            `yaml:"api-port" default:"8080"`
	Dynamic        *DynamicConfig `yaml:"-"`
}

type DynamicConfig struct {
	OutputChannels  map[string]OutputMapping `yaml:"output-channels"`
	ChannelMappings map[string][]string      `yaml:"channel-mappings"`
}

type OutputMapping struct {
	Plugin  string `yaml:"plugin"`
	Option  string `yaml:"option"`
	Channel string `yaml:"channel"`
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

	{
		d := DynamicConfig{}
		data, err = os.ReadFile(dynamicConfigPath)
		if err != nil {
			return c, err
		}
		err = yaml.Unmarshal([]byte(data), &d)
		c.Dynamic = &d
	}

	return c, err
}

func ReloadDynamicConfig(dynamicConfigPath string) (DynamicConfig, error) {
	d := DynamicConfig{}
	data, err := os.ReadFile(dynamicConfigPath)
	if err != nil {
		return d, err
	}
	err = yaml.Unmarshal([]byte(data), &d)
	if err != nil {
		return d, err
	}
	//config.Dynamic = &d
	return d, nil
}
