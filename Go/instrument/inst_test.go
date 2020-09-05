package inst_test

func main() {
	inst64 := NewInstr()
	fmt.Println("default instrument repr:", inst64)

	inst64.ID = 64
	inst64.Name = "Test_Instr_at_64"
	inst64.Address = "192.168.1.64:16064"
	inst64.State = "IN"
	fmt.Println("inst64 repr:", inst64)

	inst64UDPAddr, _ := net.ResolveUDPAddr("udp", inst64.Address)
	fmt.Printf("inst64 ip resolved: %T -> %v\n", inst64UDPAddr, inst64UDPAddr)

}