package watcher

import "context"

const (
	agentMaxCap     = 90.0 // threshold where we start doing recovery pruning
	agentTargetCap  = 80.0 // threshold we stay below in active pruning mode
	zeekPercent     = 60.0 // percentage of available storage allocated to Zeek
	suricataPercent = 30.0 // percentage of avaialble storage allocated to suricata alerts
)

var (
	storageConsumptionRate = 0.0       // Current GB/day consumption rate
	diskPercentPerDay      = 0.0       // Current disk utilization / day
	daysUntilFull          = 10        // Estimated days remaining until pruning starts
	pruneMode              = "recover" // Pruning mode, default is recover, config overrides this
)

// on agent's

// StartPruning starts data pruning tasks on the Agent.
func StartPruning(ctx context.Context) {
	// watcher always does disk pruning on the agent
	// this behavior can be disabled in the config file
	return
}

// GetAgentDataCap returns the capacity of the data volume in use on the Agent
func GetAgentDataCap() float64 {
	return 0.0
}

// GetSuriUtil returns the current utilization of the data volume in use on the Agent
func GetSuriUtil() float64 {
	// Suri is a little different since it
	// writes to a single log file and doesn't
	// do any of its own rotation.
	// We use logrotate to handle the
	// rolling and archiving of old suricata
	// alerts.

	// here we are just verifying that its
	// hapenning.  if we need to make room
	// on the disk, we'll delete archived
	// suri alert files, starting with the
	// oldest.
	return 0.0
}

// AgentAtCapacity returns true if critical disks/partitions are full
func AgentAtCapacity() {
	return
}

// PruneAgent removes old event/alert data from the Agent until it is no longer at capacity
func PruneAgent(ctx context.Context) {
	// job that gets executed regularly on the agent
	// this monitors disk usage and if a threshold is
	// crossed, begins to remove old archived zeek and
	// suricata data

	// check the mode
	// if its active, need to be below the target cap
	// if its recover, need to be below the max cap

	// if we need to, start prunining

	// find all the zeek data archives
	// sort a slice, by last modify date
	// start at the end,
	// delete directory
	// then check the current util
	// rinse/repeat until below the cap

	// update stats and checkpoints
	return
}
