package state_test

import (
	. "car3-master/Go/state"
	"fmt"
	"testing"
)

func TestStateutils(t *testing.T) {
	s0 := Undefined
	s1 := Idle
	s2 := WarmUp
	s3 := Standby
	s4 := Measure

	for _, v := range [...]State{s0, s1, s2, s3, s4} {
		fmt.Printf("state %v -> %v -> %v\n", v, StateNames[v], StateAbbr[v])
	}

	assertTrue := (s0 < s1) && (s1 < s2) && (s2 < s3) && (s3 < s4)
	if !assertTrue {
		t.Error("incorrect order of states")
	} else {
		fmt.Println("state order correct?", assertTrue)
	}

	assertTrue = StateNames[s3>>1] == "WarmUp"
	fmt.Println(StateNames[s3>>1], assertTrue)
	if !assertTrue {
		t.Error("incorrect attribution of state name")
	}

	assertTrue = StateAbbr[s4>>3] == "IN"
	fmt.Println(StateAbbr[s4>>3], assertTrue)
	if !assertTrue {
		t.Error("incorrect attribution of state name abbreviation")
	}
}