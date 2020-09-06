package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Instruments - a struct to hold them all
type Instruments struct {
	Instruments []Instrument `yaml:"Payload"`
}

// Instrument - a struct to characterize one of them
type Instrument struct {
	ID        int    `yaml:"ID"` // struct tags equal to json
	Name      string `yaml:"Name"`
	Address   string `yaml:"Address"`
	WUallowed bool   `yaml:"WU_allowed"`
}

func main() {
	// file := "D:/KIT/A350_Changeover/MasterComputer/CARIBIC3_MasterComputer/examples/config/instr_cfg.yml"
	file := "C:/Users/flori/go/src/car3-master/instrConfig/instr_cfg.yml"

	var instruments Instruments

	yamlFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err)
	}

	defer yamlFile.Close()

	yamlData, _ := ioutil.ReadAll(yamlFile)

	yaml.Unmarshal(yamlData, &instruments)

	fmt.Println(instruments)

	for i := 0; i < len(instruments.Instruments); i++ {
		fmt.Println(i)
		fmt.Printf("%+v\n", instruments.Instruments[i])
	}
}
