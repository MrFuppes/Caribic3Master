package inst_test

import (
	. "car3-master/Go/instrument"
	"car3-master/Go/state"
	"fmt"
	"net"
	"os"
	"path"
	"testing"
)

func TestInst(t *testing.T) {
	inst64 := NewInstr()
	fmt.Printf("default instrument repr:\n%v\n", inst64)

	inst64.ID = 64
	inst64.Name = "Test_Instr_at_64"
	inst64.Address = "192.168.1.64:16064"
	inst64.State, _ = state.FromAbbr("IN") // ignore error
	fmt.Printf("inst64 repr:\n%v\n", inst64)

	inst64UDPAddr, _ := net.ResolveUDPAddr("udp", inst64.Address)
	fmt.Printf("inst64 ip resolved: %T -> %v\n\n", inst64UDPAddr, inst64UDPAddr)
}

func TestPayloadFromYAML(t *testing.T) {
	wd, _ := os.Getwd()
	src := path.Join(wd, "instr_cfg_test.yml")
	p, _ := PayloadFromYAML(src)
	for k, v := range p {
		fmt.Printf("@ %v:\n%v\n", k, v)
	}
}
func TestPayloadFromJSON(t *testing.T) {
	wd, _ := os.Getwd()
	src := path.Join(wd, "instr_cfg_test.json")
	p, _ := PayloadFromJSON(src)
	for k, v := range p {
		fmt.Printf("@ %v:\n%v\n", k, v)
	}
}
