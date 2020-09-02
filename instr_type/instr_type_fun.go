//
// experimenting with a type to represent an instrument in the container
//

package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

// Instrument - a struct to characterize an instrument in the Conatiner-Lab.
type Instrument struct {
	ID        int
	Name      string
	IPAddress string
	UDPPort   int
	WUallowed bool
	State     string
}

// GetNewInstr - instantiate a new instrument
func GetNewInstr() Instrument {
	i := Instrument{}
	i.ID = -1
	i.Name = "unknown"
	i.IPAddress = "unknown"
	i.UDPPort = -1
	i.WUallowed = false
	i.State = "unknown"
	return i
}

// GetUDPAddress - method to get a UDP address string from Instrument struct.
func (i Instrument) GetUDPAddress() string {
	return i.IPAddress + ":" + strconv.Itoa(i.UDPPort)
}

// GetAddressBytes - method to get 6-byte array that represents IP address (4 bytes)
// and UPD port (2 bytes).
func (i Instrument) GetAddressBytes() [6]byte {
	result := [6]byte{}
	copy(result[0:], net.IP.To4(net.ParseIP(i.IPAddress)))
	binary.BigEndian.PutUint16(result[4:], uint16(i.UDPPort))
	return result
}

func main() {
	inst64 := GetNewInstr()
	fmt.Println("default instrument repr:", inst64)

	inst64.ID = 64
	inst64.Name = "Test_Instr_at_64"
	inst64.IPAddress = "192.168.1.64"
	inst64.UDPPort = 16064
	inst64.State = "IN"

	fmt.Println("inst64 repr:", inst64)

	// IP + port to UDP address string
	fmt.Println("inst64 IP address is:", inst64.GetUDPAddress())
	// ...and get the bytes representation:
	fmt.Println("bytes representation:", inst64.GetAddressBytes())

	// playing around with IP handlers...
	ip := net.ParseIP(inst64.IPAddress)
	b := [8]byte{1, 1, 1, 1, 1, 1, 1, 1}
	copy(b[0:], net.IP.To4(ip))
	fmt.Println("IP address in an 8 byte array:", b)

	b1 := []byte{}
	b1 = append(b1, net.IP.To4(ip)...)
	fmt.Println("IP address in a byte slice:", b1)

	// more IP parsing to and from string, to and from bytes:
	fmt.Println(net.ParseIP(inst64.IPAddress))
	fmt.Printf("%T\n", net.ParseIP(inst64.IPAddress))

	fmt.Println(net.IP.String(net.ParseIP(inst64.IPAddress)))
	fmt.Println(net.IP.To4(net.ParseIP(inst64.IPAddress)))
	fmt.Printf("%T\n", net.IP.To4(net.ParseIP(inst64.IPAddress)))
}
