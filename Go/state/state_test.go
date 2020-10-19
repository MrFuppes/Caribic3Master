package state_test

import (
	. "car3-master/state"
	"fmt"
	"testing"
)

func TestState(t *testing.T) {
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

	assertTrue = StateNames[s3-1] == "WarmUp"
	fmt.Println(StateNames[s3-1], assertTrue)
	if !assertTrue {
		t.Error("incorrect attribution of state name")
	}

	assertTrue = StateAbbr[s4-3] == "IN"
	fmt.Println(StateAbbr[s4-3], assertTrue)
	if !assertTrue {
		t.Error("incorrect attribution of state name abbreviation")
	}

}

func TestFromAbbr(t *testing.T) {
	s, e := FromAbbr("SB")
	if (s != 3) || (e != nil) {
		t.Error("abbreviation to state nbr invalid")
	}

	s, e = FromAbbr("invalid")
	if e == nil {
		t.Error("abbreviation to state nbr invalid")
	}
}

func TestFromName(t *testing.T) {
	s, e := FromName("WarmUp")
	if (s != 2) || (e != nil) {
		t.Error("name to state nbr invalid")
	}

	s, e = FromName("invalid")
	if e == nil {
		t.Error("name to state nbr invalid")
	}
}
