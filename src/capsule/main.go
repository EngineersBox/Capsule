package main

import (
	capsule "capsule/src/lib"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
	"github.com/diskfs/go-diskfs/filesystem"
)

func panicWrapper(err error) {
	panic(err)
}

var (
	props   capsule.Properties
	con     *capsule.Container
	handler capsule.Handler = capsule.Handler{
		panicWrapper,
		log.Printf,
	}
)

func main() {
	props.ReadFromJSON("container_properties.json")
	switch os.Args[1] {
	case "run":
		newUUID, err := uuid.NewRandom()
		handler.HandleErrors(err)
		con = &Container{
			newUUID,
			"container",
			RunState{false, false},
			rand.Intn(255),
			handler,
			props,
			ImageManager{
				"container",
				nil,
				filesystem.FileSystem{},
			}
		}
		con.Run(os.Args)
	case "child":
		con.SpawnChild(os.Args)
	default:
		panic("Unknown command")
	}
}
