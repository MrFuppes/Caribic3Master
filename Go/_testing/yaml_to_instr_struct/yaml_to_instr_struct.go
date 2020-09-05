package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Instruments - a struct to hold them all
type Instruments struct {
	Instruments []Instrument `yaml:"Instruments"`
}

// Instrument - a struct to characterize one of them
type Instrument struct {
	ID        int    `yaml:"ID"` // struct tags equal to json
	Name      string `yaml:"Name"`
	IPAddress string `yaml:"IP_Address"`
	UDPPort   int    `yaml:"UDP_Port"`
	WUallowed bool   `yaml:"WU_allowed"`
}

func main() {
	file := "D:/KIT/A350_Changeover/MasterComputer/CARIBIC3_MasterComputer/examples/config/instr_cfg.yml"

	var instruments Instruments

	yamlFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}

	defer yamlFile.Close()

	yamlData, _ := ioutil.ReadAll(yamlFile)

	yaml.Unmarshal(yamlData, &instruments)

	for i := 0; i < len(instruments.Instruments); i++ {
		fmt.Println(i)
		fmt.Printf("%+v\n", instruments.Instruments[i])
	}
}
