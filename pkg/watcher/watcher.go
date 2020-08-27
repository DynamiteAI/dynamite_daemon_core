package watcher

import (
	"context"
	"fmt"
	"time"

	// Dynamite
	"dynamite_daemon_core/pkg/common"
	"dynamite_daemon_core/pkg/conf"
	"dynamite_daemon_core/pkg/logging"
)

var (
	log      *logging.Entry
	logger   *logging.Logger
	exit     = make(chan []byte)
	quitting = make(chan []byte)
	running  = 0
)

// Init is called explicitly
func Init(ctx context.Context) {
	var loggers []*logging.Logger

	interval := 60 * time.Second
	worker := common.NewScheduler()
	if conf.Conf.Watcher.Interval != "" {
		watchInt, err := common.ParseDuration(conf.Conf.Watcher.Interval)
		if err == nil && watchInt != 0 {
			interval = watchInt
		} else {
			fmt.Println("Invalid Watcher interval:", err)

		}
	}

	hlog, hlogger := logging.Configure("health")
	worker.Add(ctx, WatchHealth, &hlog, interval, &quitting)
	loggers = append(loggers, &hlogger)

	if conf.Conf.HasRole("agent") {
		// run Agent monitoring tasks
		alog, alogger := logging.Configure("agent")
		worker.Add(ctx, WatchAgent, &alog, interval, &quitting)
		loggers = append(loggers, &alogger)

		// if pruning not disabled, add an Agent pruning job
		// log messages get written to health.jsonl
		// this should have its own interval...
	}
	running = worker.Count()

	go keepWatch()
	go common.Cleanup("watcher", loggers, &exit, worker)
	return
}

// keepWatch monitors the quitting channel for messages indicating a job is exiting early
// It's main job is to keep track of the running watcher jobs and signal an exit if they
// all terminate for some reason.
func keepWatch() {
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
