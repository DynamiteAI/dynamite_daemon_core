package watcher

import (
	"context"
	"strings"

	// Dynamite
	"dynamite_daemon_core/pkg/common"
	"dynamite_daemon_core/pkg/conf"
	"dynamite_daemon_core/pkg/logging"
)

var (
	fails = 0
)

// WatchAgent starts Dynamite Agent monitoring tasks
func WatchAgent(ctx context.Context, log *logging.Entry, quitting *chan []byte) {

	// start := time.Now()
	var ifaces map[string]struct{} = make(map[string]struct{})
	// setup all data collectors that should be
	// run and logged on the same interval
	suri, err := GetSuriConf()
	if err != nil {
		log.Error(err)
	} else {
		if suri.Ifaces != nil && len(suri.Ifaces) > 0 {
			for _, e := range suri.Ifaces {
				ifaces[e] = struct{}{}
			}
		}
	}

	zeek, err := GetZeekConf()
	if err != nil {
		log.Error(err)
	} else {
		if zeek.Ifaces != nil && len(zeek.Ifaces) > 0 {
			for _, e := range zeek.Ifaces {
				e = strings.Split(e, "::")[1]
				ifaces[e] = struct{}{}
			}
		}
	}

	if len(ifaces) > 0 {
		// gather packet stats for the inspection interfaces
		for k := range ifaces {
			ifrpt := GetEthInfo(k)
			if ifrpt != nil {
				rptmap := common.StructToMap(ifrpt)
				log.WithFields(rptmap).Info("interface_stats")
				fails = 0
			}
		}
	} else {
		log.Debug("Unable to retrieve inspection interfaces. Trying again.")
		fails++
	}

	if fails > 5 {
		*quitting <- []byte("Failed to identify inspection interfaces for 5 consecutive intervals.")
	}

	//
	if conf.Conf.Watcher.Mode != "" {
		rr := NewRR()
		switch conf.Conf.Watcher.Mode {
		case "engine":
			//engine
			rr.RptAEngines()
			if len(rr.EngRpts) > 0 {
				for _, v := range rr.EngRpts {
					log.WithFields(v).Info("engine_stats")
				}
			}
		case "process":
			//engine
			rr.RptAEngines()
			if len(rr.EngRpts) > 0 {
				for _, v := range rr.EngRpts {
					log.WithFields(v).Info("engine_stats")
				}
			}
			//proc
			rr.RptAProcs()
			if len(rr.ProcRpts) > 0 {
				for _, v := range rr.ProcRpts {
					log.WithFields(v).Info("proc_stats")
				}
			}

		case "thread":
			//engine
			rr.RptAEngines()
			if len(rr.EngRpts) > 0 {
				for _, v := range rr.EngRpts {
					log.WithFields(v).Info("engine_stats")
				}
			}
			//proc
			rr.RptAProcs()
			if len(rr.ProcRpts) > 0 {
				for _, v := range rr.ProcRpts {
					log.WithFields(v).Info("proc_stats")
				}
			}
			//thread
			rr.RptAThreads()
			if len(rr.ThrdRpts) > 0 {
				for _, v := range rr.ThrdRpts {
					log.WithFields(v).Info("thread_stats")
				}
			}
		}
	}

	return
}
