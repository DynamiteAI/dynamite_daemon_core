package watcher

import (
	"context"
	"fmt"
	"time"

	// Dynamite
	"github.com/DynamiteAI/dynamite_daemon_core/pkg/common"
	"github.com/DynamiteAI/dynamite_daemon_core/pkg/conf"
	"github.com/DynamiteAI/dynamite_daemon_core/pkg/logging"
)

var (
	log      *logging.Entry
	logger   *logging.Logger
	exit     = make(chan []byte)
	quitting = make(chan []byte)
	running  = 0
)

// Init is called explicitly
func Init(ctx context.Context, cfg *conf.Config) error {
	var loggers []*logging.Logger

	interval := 60 * time.Second
	worker := common.NewScheduler()
	if cfg.Watcher.Interval != "" {
		watchInt, err := common.ParseDuration(cfg.Watcher.Interval)
		if err == nil && watchInt != 0 {
			interval = watchInt
		} else {
			fmt.Println("Invalid Watcher interval:", err)

		}
	}

	hlog, hlogger, err := logging.Configure("health", cfg)
	if err != nil {
		return err
	}
	worker.Add(ctx, WatchHealth, &hlog, interval, &quitting, cfg)
	loggers = append(loggers, &hlogger)

	if cfg.HasRole("agent") {
		// run Agent monitoring tasks
		alog, alogger, err := logging.Configure("agent", cfg)
		if err != nil {
			return err
		}
		worker.Add(ctx, WatchAgent, &alog, interval, &quitting, cfg)
		loggers = append(loggers, &alogger)

		// if pruning not disabled, add an Agent pruning job
		// log messages get written to health.jsonl
		// this should have its own interval...
	}
	running = worker.Count()

	go keepWatch()
	go common.Cleanup("watcher", loggers, &exit, worker)
	return nil
}

// keepWatch monitors the quitting channel for messages indicating a job is exiting early
// It's main job is to keep track of the running watcher jobs and signal an exit if they
// all terminate for some reason.
func keepWatch() {
	for {
		select {
		case msg := <-quitting:
			logging.LogEntry.WithField("pkg", "watcher").WithField("error_msg", msg).Error("watcher_job_failed")
			running = running - 1
			if running <= 0 {
				exit <- []byte("all watcher jobs failed")
			}
		}
	}
}
