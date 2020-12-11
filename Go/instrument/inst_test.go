package inst_test

import (
	. "car3-master/instrument"
	"car3-master/state"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestInst(t *testing.T) {
	inst64 := NewInstr()
	fmt.Printf("default instrument repr:\n%v\n", inst64)

	inst64.ID = 64
	inst64.Name = "Test_Instr_at_64"
	inst64.Address = "192.168.232.64"
	inst64.State, _ = state.FromAbbr("IN") // ignore error
	fmt.Printf("inst64 repr:\n%v\n", inst64)

	inst64UDPAddr, _ := net.ResolveUDPAddr("udp", inst64.Address)
	fmt.Printf("inst64 ip resolved: %T -> %v\n\n", inst64UDPAddr, inst64UDPAddr)
}

func TestPayloadFromCfg(t *testing.T) {
	wd, _ := os.Getwd()
	// YAML
	src := path.Join(wd, "instr_cfg_test.yml")
	fmt.Printf("***\n from YAML:\n***\n")
	p, _ := PayloadFromCfg(src, yaml.Unmarshal)
	for k, v := range p {
		fmt.Printf("@ %v:\n%v\n", k, v)
	}
	// JSON
	src = path.Join(wd, "instr_cfg_test.json")
	fmt.Printf("***\n from JSON:\n***\n")
	p, _ = PayloadFromCfg(src, json.Unmarshal)
	for k, v := range p {
		fmt.Printf("@ %v:\n%v\n", k, v)
	}
}
