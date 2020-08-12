// Zeek Configuration Tools

package watcher

import (
	"strconv"
	"strings"
	"errors"
	"fmt"
	
	"gopkg.in/ini.v1"
)

type ZeekConf struct {
	MgrCPUS  []int  `json:"zeek_mgr_cpus"`
	PrxCPUS  []int  `json:"zeek_prx_cpus"`
	WrkCPUS  []int  `json:"zeek_wrk_cpus"`
	WrkProcs int    `json:"zeek_lb_procs"`
	LBMeth   string `json:"zeek_lb_method"`
	Ifaces   []string `json:"zeek_ifaces"`
}

var ZeekNodeConfFile string = "/opt/dynamite/zeek/etc/node.cfg"

func GetZeekConf() (*ZeekConf, error) {
	// Initialize a new container
	var zconf ZeekConf

	// open file as ini
	cfg, err := ini.Load(ZeekNodeConfFile)
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

			if lb_procs, ok := k["lb_procs"]; ok {
				v, err := strconv.Atoi(lb_procs)
				if err == nil {
					zconf.WrkProcs = v
				}
			}

			if lb_meth, ok := k["lb_method"]; ok {
				zconf.LBMeth = lb_meth
			}

			if iface, ok := k["interface"]; ok {
				zconf.Ifaces = append(zconf.Ifaces, iface)
			}
		}
	}

	return &zconf, nil
}
