package main

import (
	"os"

	"github.com/google/uuid"
)

var (
	props Properties = Properties{
		fsname:  "rootfs",
		procMax: "20",
		memMax:  "4096",
	}
	con *Container = nil
)

func main() {
	switch os.Args[1] {
	case "run":
		newUUID, err := uuid.NewRandom()
		handleErrors(err)
		con = &Container{
			newUUID,
			"container",
			0x00010001,
		}
		con.Run(os.Args)
	case "child":
		con.SpawnChild(os.Args)
	default:
		panic("Unknown command")
	}
}
