package conf

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	// Conf holds dynamited's current configuration settings
	Conf    Config
	roleMap = map[string]string{
		"agent":   "Inspects network traffic.",
		"monitor": "Ingests and analyzes events.",
		"scanner": "Fingerprints network devices.",
	}
)

// Config defines the structure and field names for the dynamited config file
type Config struct {
	Watcher struct {
		Interval string `yaml:"interval"`
		Mode     string `yaml:"mode"`
		Enabled  bool   `yaml:"enabled"`
	} `yaml:"watcher"`
	Scanner struct {
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

// Load initializes the conf.Conf variable with settings from the config file
func Load(s string) {
	readFile(s, &Conf)
}

func readFile(s string, cfg *Config) {
	_, err := os.Stat(s)
	if os.IsNotExist(err) {
		fmt.Printf("Config file %v not found. Exiting.\n", s)
		os.Exit(1)
	}
	f, ferr := os.Open(s)
	if ferr != nil {
		fmt.Printf("Unable to open config file %v. Exiting.\n", s)
		os.Exit(1)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	derr := decoder.Decode(cfg)
	if derr != nil {
		fmt.Printf("Unable to parse config file %v. Exiting.\n", s)
		os.Exit(1)
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
}
