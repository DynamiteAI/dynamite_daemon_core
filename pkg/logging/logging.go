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
	LogDir = ""
)

// Try to use or create the provided logging dir 
func LogDirIsUsable(s string)(bool){
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

func Init()(bool){
	// Ensure the log directory exists, try to create it if not 
	if conf.Conf.LogDir != "" {
		if LogDirIsUsable(conf.Conf.LogDir) {
			LogDir = conf.Conf.LogDir
			return true
		}
	} else {
		fmt.Println("No logging directory configured, attempting to use default /var/dynamite/managerd/logs")
		// try a default path
		if LogDirIsUsable("/var/dynamite/managerd/logs") {
			LogDir = "/var/dynamite/managerd/logs"
			return true
		} else {
			fmt.Println("Unable to configure logging. Exiting.")
			return false
		}
	}
	return false
}

func Configure(s string)(Entry, Logger){
	log := logrus.New()
	if conf.Conf.LogLevel != "" {
		if v, err := logrus.ParseLevel(conf.Conf.LogLevel); err == nil {
			log.SetLevel(v)
		} else {
			log.Error("Invalid log_level setting in config file: ", conf.Conf.LogLevel)
		}
	}
	clog := filepath.Join(LogDir,s+".jsonl")

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