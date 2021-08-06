package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	roleMap = map[string]string{
		"agent":   "Inspects network traffic.",
		"monitor": "Ingests and analyzes events.",
		"probe":   "Probes network devices for asset fingerprinting.",
	}
)

// Config defines the structure and field names for the dynamited config file
type Config struct {
	Watcher struct {
		Interval string `yaml:"interval"`
		Mode     string `yaml:"mode"`
		Enabled  bool   `yaml:"enabled"`
	} `yaml:"watcher"`
	Probe struct {
		Targets    []string `yaml:"targets"`
		Interval   string   `yaml:"interval"`
		Directory  string   `yaml:"directory"`
		Exclusions []string `yaml:"exclusions"`
		Opts       []string `yaml:"opts"`
	} `yaml:"scanner"`
	Roles    []string `yaml:"roles"`
	LogLevel string   `yaml:"log_level"`
	LogDir   string   `yaml:"log_dir"`
}

// HasRole returns true if dynamited is configured to run the given role
func (c *Config) HasRole(s string) bool {
	for _, a := range c.Roles {
		if a == s {
			return true
		}
	}
	return false
}

// Load initializes a Conf instance with settings from the config file
func Load(s string) (cfg *Config, err error) {
	cfg, err = readFile(s, cfg)
	return cfg, err
}

func readFile(s string, cfg *Config) (*Config, error) {
	_, err := os.Stat(s)
	if os.IsNotExist(err) {
		return cfg, fmt.Errorf("config file %v not found", s)
	}
	f, ferr := os.Open(s)
	if ferr != nil {
		return cfg, fmt.Errorf("unable to open config file %v", s)

	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	derr := decoder.Decode(&cfg)
	if derr != nil {
		return cfg, fmt.Errorf("unable to parse config file %v", s)
	}
	rls := []string{}
	if len(cfg.Roles) > 0 {
		for _, v := range cfg.Roles {
			if _, ok := roleMap[v]; ok {
				rls = append(rls, v)
			}
		}
		cfg.Roles = rls
	}
	return cfg, nil
}
