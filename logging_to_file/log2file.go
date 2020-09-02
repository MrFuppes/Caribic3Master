// log2file
package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// get path of current working dir:
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println("logging to ", path)

	f, err := os.OpenFile("testlogfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	// configure the logger to save UTC timestamps, display microseconds,
	// and also log the origin of the logstring
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC | log.Llongfile)
	log.SetOutput(f)

	log.Println("This is a test log entry.")
	log.Println("more logging.")
}
