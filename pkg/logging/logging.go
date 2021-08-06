package logging

import (
	// Built-ins
	"fmt"
	"os"
	"path/filepath"

	//External
	"github.com/sirupsen/logrus"

	//Dynamite
	"github.com/DynamiteAI/dynamite_daemon_core/pkg/conf"
)

type Entry struct {
	*logrus.Entry
}

type Logger struct {
	*logrus.Logger
}

var (
	// LogDir records the active logging directory
	LogDir = ""
	// DefaultLogDir is created if there is a problem creating the configured log dir
	DefaultLogDir = "/var/log/dynamite/dynamited"
	// Log can be used by other packages to write messages to dynamited log
	Log *Logger
	// LogEntry can be used by other packages to write pre-populated messages to dynamited log
	LogEntry *Entry
)

// Try to use or create the provided logging dir
func LogDirIsUsable(s string) bool {
	_, err := os.Stat(s)
	if os.IsNotExist(err) {
		err = os.MkdirAll(s, 0755)
		if err != nil {
			return false
		}
	} else if err != nil {
		return false
	}
	return true
}

func makeDefLogDir() error {
	if LogDirIsUsable(DefaultLogDir) {
		LogDir = DefaultLogDir
		return nil
	}
	return fmt.Errorf("unable to create default logging directory")
}

func createDir(cfg *conf.Config) error {
	// Ensure the log directory exists, try to create it if not
	if cfg.LogDir != "" {
		if LogDirIsUsable(cfg.LogDir) {
			LogDir = cfg.LogDir
			return nil
		} else {
			fmt.Println("Unable to use configured logging directory, attempting to use", DefaultLogDir)
			return makeDefLogDir()
		}
	}
	fmt.Println("No logging directory configured, attempting to use", DefaultLogDir)
	return makeDefLogDir()
}

func Configure(s string, cfg *conf.Config) (Entry, Logger, error) {
	log := logrus.New()
	if cfg.LogLevel != "" {
		if v, err := logrus.ParseLevel(cfg.LogLevel); err == nil {
			log.SetLevel(v)
		} else {
			return Entry{}, Logger{}, fmt.Errorf("invalid log_level setting in config file: %s", cfg.LogLevel)
		}
	}
	clog := filepath.Join(LogDir, s+".jsonl")

	LogFile, err := os.OpenFile(clog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return Entry{}, Logger{}, fmt.Errorf("unable to open %s log file, e: %v", s, err)
	}
	log.SetOutput(LogFile)
	log.SetFormatter(&logrus.JSONFormatter{})

	// Persistently set the source field
	e := log.WithFields(logrus.Fields{"source": s})

	// return the new logger and pre-populated entry
	return Entry{e}, Logger{log}, nil
}

// SetupAppLogger creates the dynamited logger and stores pointers to the logger and logger.Entry
// in global variables, LogEntry and Log. The variables can be used for threadsafe logging
// to the dynamited.jsonl log file from any other dynamited packages.
func SetupAppLogger(cfg *conf.Config) error {
	err := createDir(cfg)
	if err != nil {
		return err
	}
	le, l, err := Configure("dynamited", cfg)
	if err != nil {
		return err
	}
	LogEntry = &le
	Log = &l
	return nil
}
