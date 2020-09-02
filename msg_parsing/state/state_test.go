package stat_test

import (
	"car3-master/msg_parsing/state"
	"fmt"
	"testing"
)

func TestStateutils(t *testing.T) {
	s0 := state.Idle
	s1 := state.WarmUp

	fmt.Println(s0 < s1)
}
