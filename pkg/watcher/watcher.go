package watcher

import (
	"context"
	"time"

	// Dynamite
	"dynamite_daemon_core/pkg/logging"
	"dynamite_daemon_core/pkg/conf"
	. "dynamite_daemon_core/pkg/common"
)

var (
	log *logging.Entry
	logger *logging.Logger
	exit = make(chan []byte)
	quitting = make(chan []byte)
	running = 0
)

// Init is called explicitly
func Init(ctx context.Context) {
	log, logger := logging.Configure("watcher")
	compName := map[string]interface{}{
		"component" : "",
	}

	interval := 60 * time.Second
	worker := NewScheduler()
	if conf.Conf.Watcher.Interval != "" {
		watchInt, err := ParseDuration(conf.Conf.Watcher.Interval)
		if err == nil && watchInt != 0 {
			interval = watchInt
		} else {
			log.WithField("error_msg", err).Error("invalid_watcher_interval", err)
			
		}
	}

	if conf.Conf.HasRole("agent") {
		// run Agent monitoring tasks 
		compName["component"] = "agent"
		worker.Add(ctx, WatchAgent, &log, interval, &quitting)
		
	}
	running = worker.Count()
	
	go keepWatch()
	go Cleanup("watcher", &log, &logger, &exit, worker)
	return
}

// keepWatch monitors the quitting channel for messages indicating a job is exiting early 
// It's main job is to keep track of the running watcher jobs and signal an exit if they 
// terminate early.   
func keepWatch(){
	for {
		select {
		case msg := <-quitting:
			log.WithField("error_msg", msg).Error("watcher_job_failed")
			running = running - 1
			if running <= 0 {
				exit <- []byte("all watcher jobs failed")
			}
		}
	}
}