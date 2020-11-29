package capsule

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Properties ... Global values to constrain containerization
type Properties struct {
	ContainerName    string `json:"containerName"`
	Fsname           string `json:"fsname"`
	ProcMax          int    `json:"procMax"`
	MemMax           int    `json:"memMax"`
	TerminateOnClose bool   `json:"terminateOnClose"`
}

// ReadPropertiesFromJSON ... Parse from structured JSON file to current instance of Properties
func ReadPropertiesFromJSON(filename string) Properties {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	p := Properties{}
	log.Printf("Reading containerization properties from [%s]\n", filename)
	err = json.Unmarshal([]byte(f), &p)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Unmarshalled properties: %v\n", p)
	return p
}
