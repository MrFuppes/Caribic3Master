// Package inst implements basic types and functionality to describe and handle
// instruments (payload) of the CARIBIC container.
package inst

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Instrument - a struct to characterize an instrument in the Conatiner-Lab.
type Instrument struct {
	ID        int    `json:"ID" yaml:"ID"`
	Name      string `json:"Name" yaml:"Name"`
	Address   string `json:"Address" yaml:"Address"`
	WUallowed bool   `json:"WU_allowed" yaml:"WU_allowed"`
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

// Payload - a mapping of instrument-ID (int) -> instrument type (struct).
type Payload map[int]Instrument

type instruments struct {
	Instruments []Instrument `json:"Payload" yaml:"Payload"`
}

// PayloadFromJSON - fill the payload map with instruments from a config file in json format.
func PayloadFromJSON(src string) (Payload, error) {
	var p = Payload{}
	var inst instruments

	jsonFile, err := os.Open(src)
	if err != nil {
		return p, err
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return p, err
	}

	json.Unmarshal(jsonData, &inst)
	for i := 0; i < len(inst.Instruments); i++ {
		p[inst.Instruments[i].ID] = inst.Instruments[i]
	}
	return p, nil
}

// PayloadFromYAML - fill the payload map with instruments from a config file in yaml format.
func PayloadFromYAML(src string) (Payload, error) {
	var p = Payload{}
	var inst instruments

	yamlFile, err := os.Open(src)
	if err != nil {
		return p, err
	}
	defer yamlFile.Close()
	yamlData, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		return p, err
	}
	yaml.Unmarshal(yamlData, &inst)
	for i := 0; i < len(inst.Instruments); i++ {
		p[inst.Instruments[i].ID] = inst.Instruments[i]
	}
	return p, nil
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
