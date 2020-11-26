package capsule

import (
	"log"

	unix "golang.org/x/sys/unix"
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
	fsRoot     string
	connection unix.FdSet
	mountID    int
}

func (f *FSManager) mountFS() {
	conn, mid, err := unix.NameToFileHandle(0, f.fsRoot, 0)
	handler.HandleErrors(err)

	f.connection = conn
	f.mountID = mid

	unix.Mount(f.fsRoot, "root", "fuse", 0, "")
}
