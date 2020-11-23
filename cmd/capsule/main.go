package main

import (
	"os"

	"github.com/google/uuid"
)

var (
	cid   string     = "container"
	props Properties = Properties{"rootfs", "20"}
	con   *Container = nil
)

func main() {
	switch os.Args[1] {
	case "run":
		newUUID, err := uuid.NewRandom()
		handleErrors(err)
		con = &Container{newUUID, "container"}
		con.Run(os.Args)
	case "child":
		con.SpawnChild(os.Args)
	default:
		panic("Unknown command")
	}
}
