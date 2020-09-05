//
// experimenting with a type to represent an instrument in the container
//

package inst

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Instrument - a struct to characterize an instrument in the Conatiner-Lab.
type Instrument struct {
	ID        int
	Name      string
	Address   string
	WUallowed bool
	State     string
}

// NewInstr - instantiate a new instrument
func NewInstr() Instrument {
	i := Instrument{}
	i.ID = -1
	i.Name = "unknown"
	i.Address = "unknown"
	i.WUallowed = false
	i.State = "unknown"
	return i
}

// GetAddressBytes - method to get 6-byte array that represents IP address (4 bytes)
// and UPD port (2 bytes). Address must be of type string with format "x.x.x.x:port"
func (i Instrument) GetAddressBytes() ([6]byte, bool) {

	parts := strings.Split(i.Address, ":")
	result := [6]byte{}

	if len(parts) != 2 {
		return result, false
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return result, false
	}

	copy(result[0:], net.IP.To4(net.ParseIP(parts[0])))
	binary.BigEndian.PutUint16(result[4:], uint16(port))

	return result, true
}
