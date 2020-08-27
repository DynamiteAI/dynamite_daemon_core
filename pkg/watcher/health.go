package watcher

import (
	"context"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"

	"dynamite_daemon_core/pkg/common"
	//"dynamite_daemon_core/pkg/conf"
	"dynamite_daemon_core/pkg/logging"
)

// WatchHealth logs resource utilization stats to health.jsonl
func WatchHealth(ctx context.Context, log *logging.Entry, quitting *chan []byte) {

	// start := time.Now()
	hinfo, err := host.Info()
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("host_info_error")
	} else {
		log.WithFields(common.StructToMap(hinfo)).Info("host_info")
	}
	// CPU STATS
	ccount, err := cpu.Counts(true)
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("cpu_core_cnt_error")
	}

	lavg, err := load.Avg()
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("cpu_load_avg_error")
	} else {
		if ccount > 0 {
			log.WithFields(common.StructToMap(lavg)).WithField("cpu_cores", ccount).Info("cpu_load_avg")
		} else {
			log.WithFields(common.StructToMap(lavg)).Info("cpu_load_avg")
		}
	}

	ctimes, err := cpu.Times(true)
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("cpu_times_error")
	} else {
		for _, v := range ctimes {
			log.WithFields(common.StructToMap(v)).Info("cpu_times")
		}
	}

	// Memory
	mem, err := mem.VirtualMemory()
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("ram_stats_error")
	} else {
		log.WithFields(common.StructToMap(mem)).Info("memory_stats")
	}

	// Disk
	vpath := "/var/log/dynamite"
	vutil, err := disk.Usage(vpath)
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("disk_stats_error")
	} else {
		log.WithFields(common.StructToMap(vutil)).WithField("path", vpath).Info("disk_usage")
	}

	opath := "/opt/dynamite"
	outil, err := disk.Usage(opath)
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("disk_stats_error")
	} else {
		log.WithFields(common.StructToMap(outil)).WithField("path", opath).Info("disk_usage")
	}

	sutil, err := disk.Usage("/")
	if err != nil {
		log.WithField("error_msg", err.Error()).Error("disk_stats_error")
	} else {
		log.WithFields(common.StructToMap(sutil)).WithField("path", "/").Info("disk_usage")
	}
}
