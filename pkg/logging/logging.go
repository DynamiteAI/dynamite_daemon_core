package logging

import (
	// Built-ins
	"fmt"
	"os"
	"path/filepath"

	//External
	"github.com/sirupsen/logrus"

	//Dynamite
	"dynamite_daemon_core/pkg/conf"
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

func MakeDefLogDir() bool {
	if LogDirIsUsable(DefaultLogDir) {
		LogDir = DefaultLogDir
		return true
	} else {
		fmt.Println("Unable to create default logging directory.")
	}
	return false
}

func Init() bool {
	// Ensure the log directory exists, try to create it if not
	if conf.Conf.LogDir != "" {
		if LogDirIsUsable(conf.Conf.LogDir) {
			LogDir = conf.Conf.LogDir
			return true
		} else {
			fmt.Println("Unable to use configured logging directory, attempting to use", DefaultLogDir)
			return MakeDefLogDir()
		}
	} else {
		fmt.Println("No logging directory configured, attempting to use", DefaultLogDir)
		// try the default path
		return MakeDefLogDir()
	}
	return false
}

func Configure(s string) (Entry, Logger) {
	log := logrus.New()
	if conf.Conf.LogLevel != "" {
		if v, err := logrus.ParseLevel(conf.Conf.LogLevel); err == nil {
			log.SetLevel(v)
		} else {
			log.Error("Invalid log_level setting in config file: ", conf.Conf.LogLevel)
		}
	}
	clog := filepath.Join(LogDir, s+".jsonl")

	LogFile, err := os.OpenFile(clog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Unable to open %s log file, e: %v.\n", s, err)
		os.Exit(1)
	}
	log.SetOutput(LogFile)
	log.SetFormatter(&logrus.JSONFormatter{})

	// Persistently set the source field
	e := log.WithFields(logrus.Fields{"source": s})

	// return the new logger and pre-populated entry
	return Entry{e}, Logger{log}
}

// SetupAppLogger creates the dynamited logger and stores pointers to the logger and logger.Entry
// in global variables, LogEntry and Log. The variables can be used for threadsafe logging
// to the dynamited.jsonl log file from any other dynamited packages.
func SetupAppLogger() {
	le, l := Configure("dynamited")
	LogEntry = &le
	Log = &l
}
