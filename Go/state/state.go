// Package state implements state type for Container-Lab and instruments.
// Basic package, no dependencies (standard lib only).
package state

import (
	"errors"
	"strings"
)

// State - instrument state enum
type State uint

const (
	// Undefined state represented by uint 0
	Undefined State = iota
	// Idle state represented by uint 1
	Idle
	// WarmUp state represented by uint 2
	WarmUp
	// Standby state represented by uint 3
	Standby
	// Measure state represented by uint 4
	Measure
)

var (
	// StateNames - a mapping of State -> Name
	StateNames = map[State]string{
		0: "Undefined",
		1: "Idle",
		2: "WarmUp",
		3: "Standby",
		4: "Measure",
	}

	// StateAbbr - a mapping of State -> Abbreviation
	StateAbbr = map[State]string{
		0: "NA",
		1: "IN",
		2: "WU",
		3: "SB",
		4: "MS",
	}
)

// FromAbbr - get state number from abbreviation
func FromAbbr(expr string) (State, error) {
	expr = strings.ToUpper(expr)
	for k, v := range StateAbbr {
		if v == expr {
			return k, nil
		}
	}
	errstr := "state: invalid abbreviation -> " + expr
	return 999, errors.New(errstr)
}

// FromName - get state number from name
func FromName(expr string) (State, error) {
	expr = strings.ToLower(expr)
	for k, v := range StateNames {
		if strings.ToLower(v) == expr {
			return k, nil
		}
	}
	errstr := "state: invalid name -> " + expr
	return 999, errors.New(errstr)
}

// String converts a numeric state into its abbreviated string repr
func (s State) String() string {
	if v, ok := StateAbbr[s]; ok {
		return v
	}
	return "invalid state nbr"
}

// StringLong converts a numeric state into its verbose string repr
func (s State) StringLong() string {
	if v, ok := StateNames[s]; ok {
		return v
	}
	return "invalid state nbr"
}
