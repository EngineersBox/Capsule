package capsule

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/google/uuid"
	"golang.org/x/sys/unix"
)

const (
	mkdirPerm     os.FileMode = 0755
	writeFilePerm os.FileMode = 0700
	cgroupsDir    string      = "/sys/fs/cgroup/"
)

// RunState ... Indicates what stage of running the container is in
var RunState Enum = Enum{[]EnumItem{{0, "Stopped"}, {1, "Starting"}, {2, "Running"}}}

// ContainerState ... Object describing state of a container
type ContainerState struct {
	Running  string
	HasChild bool
}

// Container ... Containerized instance of a linux file system
type Container struct {
	ID      uuid.UUID
	Name    string
	State   ContainerState
	Cls     int
	Handler Handler
	Props   Properties
	Im      ImageManager
}

// Run ... Fork-execute an fs instance
func (c *Container) Run(args []string) {
	log.Printf("Running %v as %d\n", args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args[2:]...)...)

	c.CreateCGroup()

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Unshareflags: syscall.CLONE_NEWNS,
		Credential:   &syscall.Credential{Uid: 0, Gid: 0},
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	c.Handler.HandleErrors(cmd.Run())
	c.State.Running = RunState.Value(1)
}

// SpawnChild ... Spawn a child process within the fs instance created via run()
func (c *Container) SpawnChild(args []string) {
	log.Printf("Running in new UTS namespace %v as %d\n", args[2:], os.Getpid())

	c.Handler.HandledInvocationGroup(
		unix.Sethostname([]byte(c.Name)),
		unix.Chroot("/root/"+c.Props.Fsname),
		unix.Chdir("/"), // set the working directory inside container
		unix.Mount("proc", "proc", "proc", 0, ""),
	)

	cmd := exec.Command(args[2], args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	c.Handler.HandleErrors(cmd.Run())
	c.Handler.HandleErrors(unix.Unmount(c.Props.Fsname, 0))

	c.State.Running = RunState.Value(2)
	c.State.HasChild = true
}

// AssignCGroupMemoryAttributes ... Assign attribute values for memory usage and limiting
func (c *Container) AssignCGroupMemoryAttributes() {
	memory := filepath.Join(cgroupsDir, "memory")
	// Create memory sub-directory
	os.MkdirAll(filepath.Join(memory, c.Name), mkdirPerm)
	c.Handler.HandleErrors(
		// Set the maximum memory for the container
		ioutil.WriteFile(
			filepath.Join(memory, c.Name+"/memory.limit_in_bytes"),
			[]byte(strconv.Itoa(c.Props.MemMax)),
			writeFilePerm,
		),
	)
}

// AssignCGroupNetClsAttributes ... Assign attribute values for packet identification
func (c *Container) AssignCGroupNetClsAttributes() {
	netCls := filepath.Join(cgroupsDir, "net_cls")
	// Create net_cls sub-directory
	os.MkdirAll(filepath.Join(netCls, c.Name), mkdirPerm)
	c.Handler.HandleErrors(
		// Set the network id to identify packets from this container
		ioutil.WriteFile(
			filepath.Join(netCls, c.Name+"/net_cls.classid"),
			[]byte(strconv.Itoa(c.Cls)),
			writeFilePerm,
		),
	)
}

// CreateCGroup ... Create a CGroup for the spawned process
func (c *Container) CreateCGroup() {
	pids := filepath.Join(cgroupsDir, "pids")
	// Create CGroup sub-directory
	os.MkdirAll(filepath.Join(pids, c.Name), mkdirPerm)
	// Whether to terminate the cgroup upon all processes terminating
	willRelease := 1
	if !c.Props.TerminateOnClose {
		willRelease = 0
	}
	c.Handler.HandledInvocationGroup(
		// Set maximum child processes
		ioutil.WriteFile(
			filepath.Join(pids, c.Name+"/pids.max"),
			[]byte(strconv.Itoa(c.Props.ProcMax)),
			writeFilePerm,
		),
		// Delete the CGroup if there are no processes running
		ioutil.WriteFile(
			filepath.Join(pids, c.Name+"/notify_on_release"),
			[]byte(strconv.Itoa(willRelease)),
			writeFilePerm,
		),
		ioutil.WriteFile(
			filepath.Join(pids, c.Name+"/cgroup.procs"),
			[]byte(strconv.Itoa(os.Getpid())),
			writeFilePerm,
		),
	)
	// c.AssignCGroupMemoryAttributes()
	// c.AssignCGroupNetClsAttributes()
}

// LoadDiskImage ... Load an ISO image to create a containerized image
func (c *Container) LoadDiskImage(diskImg string) {
	c.Im.CreateIso(diskImg)
	log.Printf("Load disk image: [%s]\n", c.Im.VolumeLabel)
}
