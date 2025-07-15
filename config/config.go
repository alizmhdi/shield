package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// WhitelistSource represents a whitelist, which can be a static list of IPs, a remote URL, or both.
type WhitelistSource struct {
	IPs []string `mapstructure:"ips"`
	URL string   `mapstructure:"url"`
}

// Config is the root configuration for the shield tool.
type Config struct {
	Whitelists         map[string]WhitelistSource `mapstructure:"whitelists"`
	IngressAssignments []IngressAssignment        `mapstructure:"ingressAssignments"`
}

// IngressAssignment specifies which whitelist to apply to which ingress.
type IngressAssignment struct {
	Name      string `mapstructure:"name,omitempty"`
	Namespace string `mapstructure:"namespace,omitempty"`
	Whitelist string `mapstructure:"whitelist"`
}

// Load loads the configuration from the given file path using Viper.
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return &cfg, nil
}

// Validate checks that all assignments reference valid whitelists.
func (c *Config) Validate() error {
	for _, assign := range c.IngressAssignments {
		if _, ok := c.Whitelists[assign.Whitelist]; !ok {
			return fmt.Errorf("assignment references unknown whitelist: %s", assign.Whitelist)
		}
	}
	return nil
}
