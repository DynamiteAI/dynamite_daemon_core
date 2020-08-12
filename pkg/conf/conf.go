package conf

import (
	"gopkg.in/yaml.v2"
	"fmt"
    "os"
)

var (
	Conf Config
)

type Config struct {
    Watcher struct {
        Interval    string     `yaml:"interval"` //, envconfig:"SERVER_PORT"`
        Mode        string  `yaml:"mode"` // APP,PROCESS,THREAD,ALL
    } `yaml:"watcher"`
    Scanner struct {
        Targets []string `yaml:"targets"` //, envconfig:"DB_USERNAME"`
        Interval string `yaml:"interval"`
        Directory string `yaml:"directory"`
        Exclusions []string `yaml:"exclusions"`
        Opts []string `yaml:"opts"`
    } `yaml:"scanner"`
    Roles []string `yaml:"roles"`       //, envconfig:"DYNMGR_ROLES"
    LogLevel string `yaml:"log_level"`
    LogDir string `yaml:"log_dir"`
}

// Method for testing role membership 
func (c *Config) HasRole(s string)(bool) {
    for _, a := range c.Roles {
        if a == s {
            return true
        }
    }
    return false
}

func Load(s string)() {
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
} 
