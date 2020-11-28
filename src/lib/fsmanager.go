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

// CreateFS ... Create a filesystem at the root specified at FSManager.FsRoot
func (f *FSManager) CreateFS() {
	conn, mid, err := unix.NameToFileHandle(0, f.FsRoot, 0)
	handler.HandleErrors(err)

	f.Connection = conn
	f.MountID = mid
}

// MountFS ... Mount the filesystem specified at FSManager.FsRoot as a FUSE filesystem
func (f *FSManager) MountFS() {
	unix.Mount(f.FsRoot, "root", "fuse", 0, "")
}
