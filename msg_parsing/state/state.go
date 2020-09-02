package state

// State - instrument state enum
type State uint

const (
	Idle    State = 1 << iota // Idle / IN / [0 0 0 1]
	WarmUp                    // WarmUp / WU / [0 0 1 0]
	Standby                   // Standby / SB / [0 1 0 0]
	Measure                   // Measure / MS / [1 0 0 0]
)
