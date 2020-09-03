package watcher

import (
	"dynamite_daemon_core/pkg/common"
	"strings"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

var (
	zeekNodeTypes = map[string]struct{}{
		"manager": struct{}{},
		"proxy":   struct{}{},
		"logger":  struct{}{},
		"worker":  struct{}{},
	}

	cpuStats = map[string]struct{}{
		"user":      struct{}{},
		"system":    struct{}{},
		"idle":      struct{}{},
		"nice":      struct{}{},
		"iowait":    struct{}{},
		"irq":       struct{}{},
		"softirq":   struct{}{},
		"steal":     struct{}{},
		"guest":     struct{}{},
		"guestNice": struct{}{},
	}
)

// EngineStats is a container for resource utilization stats related to an inspection engine
type EngineStats struct {
	Engine     string  `json:"engine"`
	ProcessCnt int32   `json:"proc_cnt"`
	ThreadCnt  int32   `json:"thread_cnt"`
	CPU        float64 `json:"pct_cpu"`
	RAM        float32 `json:"pct_ram"`
	Files      int     `json:"open_files"`
}

// ProcStats is a container for resource utilization stats related to a single inspection engine process
type ProcStats struct {
	Engine   string                  `json:"engine"`
	ProcName string                  `json:"proc_name"`
	NodeID   string                  `json:"node_id"`
	PID      int32                   `json:"pid"`
	PPID     int32                   `json:"ppid"`
	CPUPct   float64                 `json:"pct_cpu"`
	RAMPct   float32                 `json:"pct_ram"`
	Files    int                     `json:"open_files"`
	Threads  int                     `json:"threads"`
	Memory   *process.MemoryInfoStat `json:"memory"`
	CPU      *cpu.TimesStat          `json:"cpu"`
}

// ThreadStats is a contianer for resource utilization stats related to a single inspection engine thread
type ThreadStats struct {
	Engine    string         `json:"engine"`
	ProcName  string         `json:"proc_name"`
	NodeID    string         `json:"node_id"`
	PID       int32          `json:"pid"`
	PPID      int32          `json:"ppid"`
	ThreadID  int32          `json:"thread_id"`
	TotalTime float64        `json:"total_time"`
	CPU       *cpu.TimesStat `json:"cpu"`
}

// AgentResourceReport is a container for managing
type AgentResourceReport struct {
	ZeekProcs []*process.Process       `json:"zeek_procs"`
	SuriProcs []*process.Process       `json:"suri_procs"`
	EngRpts   []map[string]interface{} `json:"engine_stats"`
	ProcRpts  []map[string]interface{} `json:"proc_stats"`
	ThrdRpts  []map[string]interface{} `json:"thread_stats"`
}

// NewRR generates a new AgentResourceReport instance
func NewRR() (rpt AgentResourceReport) {
	procs, err := process.Processes()
	if err != nil {
		return
	}
	for _, p := range procs {
		if n, err := p.Name(); err == nil {
			switch n {
			case "Suricata-Main":
				// found a Suri proc
				rpt.SuriProcs = append(rpt.SuriProcs, p)
			case "bro":
				// found a bro 2.x proc
				rpt.ZeekProcs = append(rpt.ZeekProcs, p)
			case "zeek":
				// found a zeek 3.x proc
				rpt.ZeekProcs = append(rpt.ZeekProcs, p)
			}
		}
	}
	return
}

// Extracts the node ID from the process command args
func getZeekNodeName(v *process.Process) string {
	if cli, err := v.CmdlineSlice(); err == nil {
		for idx, val := range cli {
			if val == "-p" {
				for t := range zeekNodeTypes {
					if strings.Contains(cli[idx+1], t) {
						return cli[idx+1]
					}
				}
			}
		}
	}
	return ""
}

// EngSummary returns aggregated stats for a given engine name
func EngSummary(procs []*process.Process, name string) map[string]interface{} {
	e := &EngineStats{Engine: name}
	for _, v := range procs {
		e.ProcessCnt++

		t, err := v.NumThreads()
		if err != nil {
			continue
		}
		e.ThreadCnt += t

		c, err := v.CPUPercent()
		if err != nil {
			continue
		}
		e.CPU += c

		m, err := v.MemoryPercent()
		if err != nil {
			continue
		}
		e.RAM += m

		f, err := v.OpenFiles()
		if err != nil {
			continue
		}
		e.Files += len(f)
	}
	return common.StructToMap(e)
}

// ProcSummary collects process-level stats for each engine
func ProcSummary(procs []*process.Process, name string) (ps []map[string]interface{}) {
	for _, proc := range procs {
		p := &ProcStats{Engine: name}

		if name == "zeek" {
			p.NodeID = getZeekNodeName(proc)
		} else if name == "suricata" {
			p.NodeID = "main"
		}

		if pn, err := proc.Name(); err == nil && pn != "" {
			p.ProcName = pn
		}

		p.PID = proc.Pid

		if pp, err := proc.Ppid(); err == nil {
			p.PPID = pp
		}

		if cp, err := proc.CPUPercent(); err == nil && cp != 0 {
			p.CPUPct = cp
		}

		if rp, err := proc.MemoryPercent(); err == nil && rp != 0 {
			p.RAMPct = rp
		}

		if of, err := proc.OpenFiles(); err == nil && of != nil {
			p.Files = len(of)
		}

		if mi, err := proc.MemoryInfo(); err == nil && mi != nil {
			p.Memory = mi
		}

		if cpu, err := proc.Times(); err == nil && cpu != nil {
			p.CPU = cpu
		}

		p.Threads = 0
		if thrds, err := proc.Threads(); err == nil && len(thrds) > 0 {
			p.Threads = len(thrds)
		}

		ps = append(ps, common.StructToMap(p))
	}
	return
}

// ThreadSummary collects thread-level stats for each engine, process
func ThreadSummary(procs []*process.Process, name string) (ts []map[string]interface{}) {
	for _, proc := range procs {

		if thrds, err := proc.Threads(); err == nil && len(thrds) > 0 {
			for id, thrd := range thrds {
				p := &ThreadStats{Engine: name, ThreadID: id}

				if name == "zeek" {
					p.NodeID = getZeekNodeName(proc)
				} else if name == "suricata" {
					p.NodeID = "main"
				}

				if pn, err := proc.Name(); err == nil && pn != "" {
					p.ProcName = pn
				}

				p.PID = proc.Pid

				if pp, err := proc.Ppid(); err == nil {
					p.PPID = pp
				}
				p.CPU = thrd
				ts = append(ts, common.StructToMap(p))
			}
		}
	}
	return
}

// RptAEngines stores a slice of engine stat reports in ResourceReport.EngRpts
func (rpt *AgentResourceReport) RptAEngines() {
	if rpt.ZeekProcs != nil && len(rpt.ZeekProcs) > 0 {
		rpt.EngRpts = append(rpt.EngRpts, EngSummary(rpt.ZeekProcs, "zeek"))
	}
	if rpt.SuriProcs != nil && len(rpt.SuriProcs) > 0 {
		rpt.EngRpts = append(rpt.EngRpts, EngSummary(rpt.SuriProcs, "suricata"))
	}
	return
}

// RptAProcs stores a slice of process stat reports in ResourceReport.ProcRpts
func (rpt *AgentResourceReport) RptAProcs() {
	if rpt.ZeekProcs != nil && len(rpt.ZeekProcs) > 0 {
		rpt.ProcRpts = append(rpt.ProcRpts, ProcSummary(rpt.ZeekProcs, "zeek")...)
	}
	if rpt.SuriProcs != nil && len(rpt.SuriProcs) > 0 {
		rpt.ProcRpts = append(rpt.ProcRpts, ProcSummary(rpt.SuriProcs, "suricata")...)
	}
	return
}

// RptAThreads stores a slice of thread stat reports in ResourceReport.ThrdRpts
func (rpt *AgentResourceReport) RptAThreads() {
	if rpt.ZeekProcs != nil && len(rpt.ZeekProcs) > 0 {
		for _, v := range ThreadSummary(rpt.ZeekProcs, "zeek") {
			if len(v) > 0 {
				rpt.ThrdRpts = append(rpt.ThrdRpts, v)
			}
		}
	}
	if rpt.SuriProcs != nil && len(rpt.SuriProcs) > 0 {
		for _, v := range ThreadSummary(rpt.SuriProcs, "suricata") {
			if len(v) > 0 {
				rpt.ThrdRpts = append(rpt.ThrdRpts, v)
			}
		}
	}
	return
}
