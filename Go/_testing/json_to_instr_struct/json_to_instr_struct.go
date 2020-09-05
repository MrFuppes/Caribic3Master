package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Instruments - a struct to hold them all
type Instruments struct {
	Instruments []Instrument `json:"Instruments"`
}

// Instrument - a struct to characterize one of them
type Instrument struct {
	ID        int    `json:"ID"`
	Name      string `json:"Name"`
	IPAddress string `json:"IP_Address"`
	UDPPort   int    `json:"UDP_Port"`
	WUallowed bool   `json:"WU_allowed"`
}

func main() {
	file := "D:/KIT/A350_Changeover/MasterComputer/CARIBIC3_MasterComputer/examples/config/instr_cfg.json"

	var allInstruments Instruments

	jsonFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	jsonData, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(jsonData, &allInstruments)

	for i := 0; i < len(allInstruments.Instruments); i++ {
		fmt.Println(i)
		fmt.Printf("%+v\n", allInstruments.Instruments[i])

	}
}
