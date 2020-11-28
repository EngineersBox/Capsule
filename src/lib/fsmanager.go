package capsule

import (
	"log"

	"golang.org/x/sys/unix"
)

func panicWrapper(err error) {
	panic(err)
}

var handler Handler = Handler{
	panicWrapper,
	log.Printf,
}

// FSManager ... A manager to handle filesystem operations
type FSManager struct {
	FsRoot     string
	Connection unix.FdSet
	MountID    int
}

func (f *FSManager) mountFS() {
	unix.Mount(f.FsRoot, "root", "fuse", 0, "")
}
