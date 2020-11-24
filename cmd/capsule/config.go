package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Properties ... Global values to constrain containerization
type Properties struct {
	fsname  string
	procMax string
	memMax  string
}

func (p *Properties) readFromJSON(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reading containerization properties from [%s]\n", filename)
	defer f.Close()
	propBytes, _ := ioutil.ReadAll(f)
	json.Unmarshal(propBytes, p)
}
