package main

import (
	capsule "capsule/src/lib"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
)

func panicWrapper(err error) {
	panic(err)
}

var (
	props   capsule.Properties
	con     *capsule.Container
	invoker capsule.Invoker = capsule.Invoker{
		panicWrapper,
		log.Printf,
	}
)

func main() {
	props.ReadFromJSON("container_properties.json")
	switch os.Args[1] {
	case "run":
		newUUID, err := uuid.NewRandom()
		invoker.HandleErrors(err)
		con = &Container{
			newUUID,
			"container",
			RunState{false, false},
			rand.Intn(255),
			invoker,
			props,
		}
		con.Run(os.Args)
	case "child":
		con.SpawnChild(os.Args)
	default:
		panic("Unknown command")
	}
}
