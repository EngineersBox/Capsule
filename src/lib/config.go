package capsule

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Properties ... Global values to constrain containerization
type Properties struct {
	fsname  string
	procMax int
	memMax  int
}

// ReadFromJSON ... Parse from structured JSON file to current instance of Properties
func (p *Properties) ReadFromJSON(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reading containerization properties from [%s]\n", filename)
	defer f.Close()
	propBytes, _ := ioutil.ReadAll(f)
	json.Unmarshal(propBytes, p)
}
