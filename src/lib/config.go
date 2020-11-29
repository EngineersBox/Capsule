package capsule

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Properties ... Global values to constrain containerization
type Properties struct {
	Fsname  string `json:"fsname"`
	ProcMax string `json:"procMax"`
	MemMax  string `json:"memMax"`
}

// ReadFromJSON ... Parse from structured JSON file to current instance of Properties
func (p *Properties) ReadFromJSON(filename string) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Reading containerization properties from [%s]\n", filename)
	err = json.Unmarshal([]byte(f), p)
	if err != nil {
		log.Fatal(err)
	}
}
