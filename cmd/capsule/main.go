package main

import (
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
)

func panicWrapper(err error) {
	panic(err)
}

var (
	props   Properties
	con     *Container
	invoker Invoker = Invoker{
		panicWrapper,
		log.Printf,
	}
)

func main() {
	props.readFromJSON("container_properties.json")
	switch os.Args[1] {
	case "run":
		newUUID, err := uuid.NewRandom()
		invoker.handleErrors(err)
		con = &Container{
			newUUID,
			"container",
			rand.Intn(255),
			invoker,
		}
		con.Run(os.Args)
	case "child":
		con.SpawnChild(os.Args)
	default:
		panic("Unknown command")
	}
}
