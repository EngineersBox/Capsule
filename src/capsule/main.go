package main

import (
	capsule "capsule/src/lib"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
)

var (
	props   capsule.Properties
	con     *capsule.Container
	handler capsule.Handler = capsule.Handler{
		ErrorHandler: func(err error) {
			panic(err)
		},
		InfoHandler: log.Printf,
	}
)

func main() {
	props.ReadFromJSON("config/container_properties.json")
	switch os.Args[1] {
	case "run":
		newUUID, err := uuid.NewRandom()
		handler.HandleErrors(err)
		con = &capsule.Container{
			ID:   newUUID,
			Name: "container",
			State: capsule.RunState{
				Running:  false,
				HasChild: false,
			},
			Cls:     rand.Intn(255),
			Handler: handler,
			Props:   props,
			Im: capsule.ImageManager{
				VolumeLabel: "container",
				FsDisk:      nil,
				Fs:          nil,
			},
		}
		con.Run(os.Args)
	case "child":
		con.SpawnChild(os.Args)
	default:
		panic("Unknown command")
	}
}
