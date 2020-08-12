package watcher

// monitor the Agent's health and performance

import (
	"net"
	"os"
	"os/exec"
	"fmt"
	"strings"
	"errors"
	"path/filepath"
	"strconv"
	"github.com/safchain/ethtool"
	"dynamite_daemon_core/pkg/common"
)

var (
	stats = map[string]struct{}{
		"collisions": 		struct{}{},
		"multicast":        struct{}{},
		"rx_bytes":         struct{}{},
		"rx_compressed":    struct{}{},
		"rx_crc_errors":    struct{}{},
		"rx_dropped":       struct{}{},
		"rx_errors":        struct{}{},
		"rx_fifo_errors":   struct{}{},
		"rx_frame_errors":  struct{}{},
		"rx_length_errors": struct{}{},
		"rx_missed_errors": struct{}{},
		"rx_nohandler":     struct{}{},
		"rx_over_errors":   struct{}{},
		"rx_packets":       struct{}{}}
)

type Iface struct {
	Name			string	`json:"name"`
	Driver			string	`json:"driver"`
	DrvVer			string	`json:"driver_ver"`
	Link			bool	`json:"link"`
	Collisions		int64  `json:"collisions"`
	Mulitcast		int64  `json:"multicast"`
	Bytes			int64  `json:"bytes"`
	Compressed		int64  `json:"compressed"`
	CRC_errors		int64  `json:"crc_errors"`
	Dropped			int64  `json:"dropped"`
	Errors			int64  `json:"errors"`
	FIFO_errors		int64  `json:"fifo_errors"`
	Frame_errors	int64  `json:"frame_errors"`
	Length_errors	int64  `json:"length_errors"`
	Missed_errors	int64  `json:"missed_errors"`
	NOHandler		int64  `json:"nohandler"`
	Over_errors		int64  `json:"over_errors"`
	Packets			int64  `json:"packets"`
}

// checks if a given interface is valid
func isIface(i string) bool {
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, v := range ifaces {
			if i == v.Name {
				return true
			}
		}
	} else {
		fmt.Println("Interface validation error:", err)
	}
	return false
}

// function to validate given interface name and return 
// a prepopulated Iface struct 
func newIface(s string)(iface *Iface, err error) {
	if isIface(s) {
		iface = &Iface{Name:s}
		return iface, nil
	}
	return nil, errors.New("invalid interface name")
}

// Method to update Iface stats using sysfs 
func (iface *Iface) updateStats()(){
	stats := readSysFS(iface.Name)
	if stats != nil {
		for k,v := range stats {
			switch k {
			case "collisions":
				iface.Collisions = v
				break
			case "multicast":
				iface.Mulitcast = v
				break
			case "rx_bytes":
				iface.Bytes = v
				break
			case "rx_compressed":
				break
				iface.Compressed = v
				break
			case "rx_crc_errors":
				iface.CRC_errors = v
				break
			case "rx_dropped":
				iface.Dropped = v
				break
			case "rx_errors":
				iface.Errors = v
				break
			case "rx_fifo_errors":
				iface.FIFO_errors = v
				break
			case "rx_frame_errors":
				iface.Frame_errors = v
				break
			case "rx_length_errors":
				iface.Length_errors = v
				break
			case "rx_missed_errors":
				iface.Missed_errors = v
				break
			case "rx_nohandler":
				iface.NOHandler = v
				break
			case "rx_over_errors":
				iface.Over_errors = v
				break
			case "rx_packets":
				iface.Packets = v
			}
		}
	}
	return
}

// takes an interface name and returns the available sysfs counters
func readSysFS(ifname string) (ifrpt map[string]int64) {
	ifrpt = make(map[string]int64)
	if isIface(ifname) {
		statspath := filepath.Join("/sys/class/net", ifname, "statistics")
		avail := common.GetFileList(statspath)
		if avail != nil {
			for _, v := range avail {
				if _, ok := stats[v]; ok {
					statpath := filepath.Join(statspath, v)
					val := runCMD("cat", []string{statpath})
					if val != "" {
						ival, err := strconv.ParseInt(strings.TrimSpace(val), 0, 64)
						if err != nil {
							fmt.Printf("Unable to parse stat %s: %v\n", v, err)
						} else {
							ifrpt[v] = ival
						}
					}
				}
			}
		}
	}
	return
}

// get interface link state 
func (iface *Iface)getLinkStatus(e *ethtool.Ethtool)() {
	ls, err := e.LinkState(iface.Name)
	if err != nil {
		iface.Link = false
		return
	}
	if ls == 1 {
		iface.Link = true
	} else {
		iface.Link = false
	}
	return 
}

// function to gather network interface info and stats 
func GetEthInfo(i string)(iface *Iface) {
	iface, err := newIface(i)
	if err != nil {
		return 
	}
	et, err := ethtool.NewEthtool()
	if err != nil {
		return
	}
	defer et.Close()
	iface.getLinkStatus(et)
	iface.getDriverInfo(et)
	iface.updateStats()
	return 
}

// get inteface driver name 
func (iface *Iface)getDriverInfo(e *ethtool.Ethtool)() {
	dinfo, err := e.DriverInfo(iface.Name)
	if err != nil {
		return 
	}
	if dinfo.Driver != "" {
		iface.Driver = dinfo.Driver
	}
	if dinfo.Version != "" {
		iface.DrvVer = dinfo.Version
	}
	return 
}

// returns a list of files found in the given directory
func getStatList(dir string) (files []string) {
	f, err := os.Open(dir)
	if err != nil {
		return nil
	}
	dirlist, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil
	}
	for _, e := range dirlist {
		files = append(files, e.Name())
	}
	return
}

// execute a given command, passing arguments.
//returns stdout as a string or nil
func runCMD(c string, args []string) (out string) {
	cmd := exec.Command(c, args[:]...)
	res, err := cmd.Output()
	if err != nil {
		out = ""
	} else {
		out = string(res)
	}
	return
}