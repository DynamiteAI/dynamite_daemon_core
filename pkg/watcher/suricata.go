// Suricata Configuration Tools

package watcher

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"reflect"
	"errors"
	"fmt"
)

// Struct to hold Suricata threading settings
type SuriConf struct {
	MgrCPUS     	[]int   	`json:"suri_mgr_cpus"`
	RcvCPUS     	[]int   	`json:"suri_rcv_cpus"`
	WrkCPUS     	[]int   	`json:"suri_wrk_cpus"`
	DTR         	float64 	`json:"suri_dtr"`
	UseAffinity 	bool    	`json:"suri_pin_cpus"`
	RunMode     	string  	`json:"suri_runmode"`
	Ifaces			[]string 	`json:"suri_ifaces"`
}

// Variables used in this package
var (
	SuriConfFile = "/etc/dynamite/suricata/suricata.yaml"
)

// Function to parse threading settings from Suricata yaml
func GetSuriConf() (*SuriConf, error) {
	var sconf SuriConf
	if _, err := os.Stat(SuriConfFile); os.IsNotExist(err) {
		msg := fmt.Sprintf("Suricata conf file not found, e: %v", err)
		return &sconf, errors.New(msg)
	}

	scfile, err := ioutil.ReadFile(SuriConfFile)
	if err != nil {
		msg := fmt.Sprintf("Failed to open Suricata config file, e: %v")
		return &sconf, errors.New(msg)
	}

	sc := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(scfile), &sc); err != nil {
		msg := fmt.Sprintf("Error parsing Suri conf: %v", err)
		return &sconf, errors.New(msg)
	}

	if rm, ok := sc["runmode"]; ok && rm != nil {
		sconf.RunMode = rm.(string)
	}

	if af, ok := sc["af-packet"]; ok && af != nil {
		for _, v := range af.([]interface{}) {
			for vk, vv := range v.(map[interface{}]interface{}) {
				if reflect.TypeOf(vk).String() == "string" && vk.(string) == "interface" {
					sconf.Ifaces = append(sconf.Ifaces, vv.(string))
				}
			}
		}
	}

	if thd, ok := sc["threading"]; ok && thd != nil {
		for k, v := range thd.(map[interface{}]interface{}) {
			if k.(string) == "set-cpu-affinity" {
				sconf.UseAffinity = v.(bool)
			}

			if k.(string) == "detect-thread-ratio" {
				sconf.DTR = v.(float64)
			}

			if k.(string) == "cpu-affinity" {
				cpua := v.(interface{})
				for i := range cpua.([]interface{}) {
					imap := cpua.([]interface{})[i].(map[interface{}]interface{})
					for k, v := range imap {
						if k.(string) == "management-cpu-set" {
							for k, v := range v.(map[interface{}]interface{}) {
								if k.(string) == "cpu" {
									if reflect.TypeOf(v).String() == "[]interface {}" {
										for c := range v.([]interface{}) {
											sconf.MgrCPUS = append(sconf.MgrCPUS, c)									
										}
									} else if reflect.TypeOf(v).String() == "int" {
										sconf.MgrCPUS = append(sconf.MgrCPUS, v.(int))
									}
								}
							}
						}

						if k.(string) == "receive-cpu-set" {
							for k, v := range v.(map[interface{}]interface{}) {
								if k.(string) == "cpu" {
									if reflect.TypeOf(v).String() == "[]interface {}" {
										for c := range v.([]interface{}) {
											sconf.RcvCPUS = append(sconf.RcvCPUS, c)
											}
										}
									} else if reflect.TypeOf(v).String() == "int" {
										sconf.RcvCPUS = append(sconf.RcvCPUS, v.(int))
									}
								}
							}
						if k.(string) == "worker-cpu-set" {
							for k, v := range v.(map[interface{}]interface{}) {
								if k.(string) == "cpu" {
									if reflect.TypeOf(v).String() == "[]interface {}" {
										for c := range v.([]interface{}) {
											sconf.WrkCPUS = append(sconf.WrkCPUS, c)
											}
									} else if reflect.TypeOf(v).String() == "int" {
										sconf.WrkCPUS = append(sconf.WrkCPUS, v.(int))
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return &sconf, nil
}
