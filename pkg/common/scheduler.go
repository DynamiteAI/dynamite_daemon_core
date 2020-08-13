package common

import (
	"context"
	"sync"
	"time"
	"dynamite_daemon_core/pkg/logging"
	"os"
	"fmt"
)

type Job func(ctx context.Context, log *logging.Entry, exit *chan []byte)

type Scheduler struct {
	wg            *sync.WaitGroup
	cancellations []context.CancelFunc
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		wg:            new(sync.WaitGroup),
		cancellations: make([]context.CancelFunc, 0),
	}
}

// Count returns the number of scheduled jobs 
func (s *Scheduler) Count()(int) {
	return len(s.cancellations)
}

// Add starts goroutine which constantly calls provided job with interval delay
func (s *Scheduler) Add(ctx context.Context, j Job, log *logging.Entry, interval time.Duration, quitting *chan []byte) {
	ctx, cancel := context.WithCancel(ctx)
	s.cancellations = append(s.cancellations, cancel)

	s.wg.Add(1)
	go s.process(ctx, j, log, interval, quitting)
}

// Stop cancels all running jobs
func (s *Scheduler) Stop() {
	for _, cancel := range s.cancellations {
		cancel()
	}
	s.wg.Wait()
}

func (s *Scheduler) process(ctx context.Context, j Job, log *logging.Entry, interval time.Duration, quitting *chan []byte) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			j(ctx, log, quitting)
		case <-ctx.Done():
			s.wg.Done()
			return
		}
	}
}

// Function used by dynamited packages to defer cleanup of their things 
func Cleanup(pkg string, loggers []*logging.Logger, exit *chan []byte, worker *Scheduler){
	for {
		select {
		case msg := <-*exit:
			fmt.Printf("%s is exiting. Error message: %s", pkg, string(msg))
			for _, v := range loggers {
				if file, ok := v.Out.(*os.File); ok {
					file.Sync()
					file.Close()
				}
			}
			worker.Stop() 
		}
	}
}