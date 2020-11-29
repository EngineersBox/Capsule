package main

import (
	capsule "capsule/src/lib"
	"log"
	"math/rand"
	"os"

	"github.com/google/uuid"
)

func flattenArgs(a ...interface{}) []interface{} {
	return a
}

var (
	props capsule.Properties = capsule.ReadPropertiesFromJSON("config/container_properties.json")
	con   capsule.Container  = capsule.Container{
		ID:   flattenArgs(uuid.NewRandom())[0].(uuid.UUID),
		Name: props.ContainerName,
		State: capsule.RunState{
			Running:  false,
			HasChild: false,
		},
		Cls:     rand.Intn(255),
		Handler: handler,
		Props:   props,
		Im: capsule.ImageManager{
			VolumeLabel: props.ContainerName,
			FsDisk:      nil,
			Fs:          nil,
		},
	}
	handler capsule.Handler = capsule.Handler{
		ErrorHandler: func(err error) {
			panic(err)
		},
		InfoHandler: log.Printf,
	}
)

func main() {
	switch os.Args[1] {
	case "run":
		con.Run(os.Args)
	case "child":
		con.SpawnChild(os.Args)
	default:
		panic("Unknown command")
	}
}
