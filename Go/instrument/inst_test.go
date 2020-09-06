package inst_test

import (
	. "car3-master/Go/instrument"
	"fmt"
	"net"
	"os"
	"path"
	"testing"
)

func TestInst(t *testing.T) {
	inst64 := NewInstr()
	fmt.Println("default instrument repr:", inst64)
	b, ok := inst64.GetAddressBytes()
	fmt.Println(b, ok)

	inst64.ID = 64
	inst64.Name = "Test_Instr_at_64"
	inst64.Address = "192.168.1.64:16064"
	inst64.State = "IN"
	fmt.Println("inst64 repr:", inst64)
	inst64UDPAddr, _ := net.ResolveUDPAddr("udp", inst64.Address)
	fmt.Printf("inst64 ip resolved: %T -> %v\n", inst64UDPAddr, inst64UDPAddr)

	b, ok = inst64.GetAddressBytes()
	fmt.Println(b, ok)
}

func TestPayloadFromYAML(t *testing.T) {
	wd, _ := os.Getwd()
	src := path.Join(wd, "instr_cfg_test.yml")
	p, _ := PayloadFromYAML(src)
	for k, v := range p {
		fmt.Println(k, v)
	}
}
func TestPayloadFromJSON(t *testing.T) {
	wd, _ := os.Getwd()
	src := path.Join(wd, "instr_cfg_test.json")
	p, _ := PayloadFromJSON(src)
	for k, v := range p {
		fmt.Println(k, v)
	}
}
