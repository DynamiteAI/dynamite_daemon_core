package watcher

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/host"

	. "dynamite_daemon_core/pkg/common"
	//"dynamite_daemon_core/pkg/conf"
	"dynamite_daemon_core/pkg/logging"
)

// 
func WatchHealth(ctx context.Context, log *logging.Entry, quitting *chan []byte) {

	// start := time.Now()
	hinfo, err := host.Info()
	if err != nil {
		fmt.Println("Unable to retrieve host information.")
		*quitting <- []byte("Health watcher exiting!!")
		return
	} else {
		log.WithFields(StructToMap(hinfo)).Info("host_info")	
	}

	// CPU STATS
	ccount, err:= cpu.Counts(true)
	if err != nil {
		fmt.Println("Unable to retrieve the cpu core count.")
		*quitting <- []byte("Health watcher exiting!!")		
		return 
	} 

	lavg, err := load.Avg()
	if err != nil {
		fmt.Println("Unable to retrieve the cpu load average.")
		*quitting <- []byte("Health watcher exiting!!")
		return
	} else {
		if ccount > 0 {
			log.WithFields(StructToMap(lavg)).WithField("cpu_cores", ccount).Info("cpu_load_avg")	
		} else {
			log.WithFields(StructToMap(lavg)).Info("cpu_load_avg")
		}
	}

	ctimes, err := cpu.Times(true)
	if err != nil {
		fmt.Println("Unable to retrieve the cpu times.")
		*quitting <- []byte("Health watcher exiting!!")
		return 
	} else {
		for _, v := range ctimes {
			log.WithFields(StructToMap(v)).Info("cpu_times")
		}
	}
		
	// Memory 
	mem, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Unable to retrieve RAM stats.")
		*quitting <- []byte("Health watcher exiting!!")
		return 
	} else {
		log.WithFields(StructToMap(mem)).Info("memory_stats")
	}

	// Disk
	vutil, err := disk.Usage("/var/dynamite")
	if err != nil {
		fmt.Println("Unable to retrieve disk usage stats.")
		*quitting <- []byte("Health watcher exiting!!")
		return 		
	} else { 
		log.WithFields(StructToMap(vutil)).WithField("path", "/var/dynamite").Info("disk_usage")
	}

	outil, err := disk.Usage("/opt/dynamite")
	if err != nil {
		fmt.Println("Unable to retrieve disk usage stats.")
		*quitting <- []byte("Health watcher exiting!!")
		return 		
	} else { 
		log.WithFields(StructToMap(outil)).WithField("path", "/opt/dynamite").Info("disk_usage")
	}

	sutil, err := disk.Usage("/")
	if err != nil {
		fmt.Println("Unable to retrieve disk usage stats.")
		*quitting <- []byte("Health watcher exiting!!")
		return 		
	} else { 
		log.WithFields(StructToMap(sutil)).WithField("path", "/").Info("disk_usage")
	}
}