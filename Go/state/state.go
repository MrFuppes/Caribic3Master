package state

// State - instrument state enum
type State uint

// Idle / IN / [0 0 0 1]
// WarmUp / WU / [0 0 1 0]
// Standby / SB / [0 1 0 0]
// Measure / MS / [1 0 0 0]
const (
	Undefined State = 0
	Idle      State = 1
	WarmUp    State = 2
	Standby   State = 4
	Measure   State = 8
)

// StateNames - a mapping of State -> Name
var StateNames = map[State]string{
	0: "Undefined",
	1: "Idle",
	2: "WarmUp",
	4: "Standby",
	8: "Measure",
}

// StateAbbr - a mapping of State -> Abbreviation
var StateAbbr = map[State]string{
	0: "NA",
	1: "IN",
	2: "WU",
	4: "SB",
	8: "MS",
}
