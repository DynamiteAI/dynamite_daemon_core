// Zeek Configuration Tools

package watcher

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

// ZeekConf is a struct that contains Zeek's network interface and processor settings
type ZeekConf struct {
	MgrCPUS  []int    `json:"zeek_mgr_cpus"`
	PrxCPUS  []int    `json:"zeek_prx_cpus"`
	WrkCPUS  []int    `json:"zeek_wrk_cpus"`
	WrkProcs int      `json:"zeek_lb_procs"`
	LBMeth   string   `json:"zeek_lb_method"`
	Ifaces   []string `json:"zeek_ifaces"`
}

var zeekNodeConfig string = "/opt/dynamite/zeek/etc/node.cfg"

// GetZeekConf parses the Zeek node config and returns a populated ZeekConf struct
func GetZeekConf() (*ZeekConf, error) {
	// Initialize a new container
	var zconf ZeekConf

	// open file as ini
	cfg, err := ini.Load(zeekNodeConfig)
	if err != nil {
		msg := fmt.Sprintf("Fail to read Zeek conf file: %v", err)
		return &zconf, errors.New(msg)
	}

	// get list of sections
	sections := cfg.Sections()
	// for each section, see if name contains string worker

	for s := range sections {
		if strings.Contains(sections[s].Name(), "worker") {
			k := sections[s].KeysHash()
			if k == nil || len(k) == 0 {
				continue
			}
			if wcpus, ok := k["pin_cpus"]; ok {
				for p := range strings.Split(wcpus, ",") {
					zconf.WrkCPUS = append(zconf.WrkCPUS, p)
				}
			}

			if LBProcs, ok := k["lb_procs"]; ok {
				v, err := strconv.Atoi(LBProcs)
				if err == nil {
					zconf.WrkProcs = v
				}
			}

			if LBMethod, ok := k["lb_method"]; ok {
				zconf.LBMeth = LBMethod
			}

			if iface, ok := k["interface"]; ok {
				zconf.Ifaces = append(zconf.Ifaces, iface)
			}
		}
	}

	return &zconf, nil
}
