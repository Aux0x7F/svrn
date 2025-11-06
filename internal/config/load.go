package config

import (
	"errors"
	"flag"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is the fully-resolved runtime configuration after merging
// defaults, YAML file, env vars, and flags.
type Config struct {
	Roles     []string `yaml:"roles"`
	Services  []string `yaml:"services"`
	Community string   `yaml:"community"`
	Router    string   `yaml:"router"`
}

// Load resolves configuration from flags, env, and optional YAML file.
func Load() (*Config, error) {
	flagRoles := flag.String("roles", "", "comma-separated: consumer,provider,relay,seed (default consumer-only)")
	flagServices := flag.String("services", "", "comma-separated: blob,crdt")
	flagConfig := flag.String("config", "", "path to YAML config (optional)")
	flagCommunity := flag.String("community", "", "bootstrap community URI")
	flagRouter := flag.String("router", "auto", "i2p router: auto | external:host:port")

	// Prevent flag.Parse() from running twice
	if !flag.Parsed() {
		flag.Parse()
	}

	cfg := &Config{
		Roles:     nil,
		Services:  nil,
		Community: os.Getenv("SVRN_COMMUNITY"),
		Router:    envOrDefault("SVRN_ROUTER", "auto"),
	}

	// Load YAML config if provided
	if *flagConfig != "" {
		data, err := os.ReadFile(*flagConfig)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// Apply env overrides
	if v := os.Getenv("SVRN_ROLES"); v != "" {
		cfg.Roles = splitCSV(v)
	}
	if v := os.Getenv("SVRN_SERVICES"); v != "" {
		cfg.Services = splitCSV(v)
	}

	// Apply flag overrides (highest precedence)
	if *flagRoles != "" {
		cfg.Roles = splitCSV(*flagRoles)
	}
	if *flagServices != "" {
		cfg.Services = splitCSV(*flagServices)
	}
	if *flagCommunity != "" {
		cfg.Community = *flagCommunity
	}
	if *flagRouter != "" {
		cfg.Router = *flagRouter
	}

	// Default role if none specified
	if len(cfg.Roles) == 0 {
		cfg.Roles = []string{"consumer"}
	}

	return cfg, nil
}

func splitCSV(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Validate ensures config is internally consistent.
func (c *Config) Validate() error {
	validRoles := map[string]bool{"consumer": true, "provider": true, "relay": true, "seed": true}
	for _, r := range c.Roles {
		if !validRoles[r] {
			return errors.New("invalid role: " + r)
		}
	}
	return nil
}
