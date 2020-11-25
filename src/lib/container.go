package capsule

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"

	"github.com/google/uuid"
)

const (
	mkdirPerm     os.FileMode = 0755
	writeFilePerm os.FileMode = 0700
)

// RunState ... Object describing state of a container
type RunState struct {
	running  bool
	hasChild bool
}

// Container ... Containerized instance of a linux file system
type Container struct {
	id      uuid.UUID
	name    string
	state   RunState
	cls     int
	invoker Invoker
	props   Properties
	m       sync.Mutex
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
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	c.invoker.HandleErrors(cmd.Run())
	c.state.running = true
}

// SpawnChild ... Spawn a child process within the fs instance created via run()
func (c *Container) SpawnChild(args []string) {
	log.Printf("Running in new UTS namespace %v as %d\n", args[2:], os.Getpid())

	c.invoker.HandledInvocationGroup(
		syscall.Sethostname([]byte(c.name)),
		syscall.Chroot("/root/"+c.props.fsname),
		syscall.Chdir("/"), // set the working directory inside container
		syscall.Mount("proc", "proc", "proc", 0, ""),
	)

	cmd := exec.Command(args[2], args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	c.invoker.HandleErrors(cmd.Run())
	c.invoker.HandleErrors(syscall.Unmount(c.props.fsname, 0))

	c.state.hasChild = true
}

// CreateCGroup ... Create a CGroup for the spawned process
func (c *Container) CreateCGroup() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	memory := filepath.Join(cgroups, "memory")
	netCls := filepath.Join(cgroups, "net_cls")

	c.invoker.HandledInvocationGroup(
		// Create CGroup sub-directory
		os.Mkdir(filepath.Join(pids, c.name), mkdirPerm),
		// Set maximum child processes
		ioutil.WriteFile(filepath.Join(pids, c.name+"/pids.max"), []byte(c.props.procMax), writeFilePerm),
		// Delete the CGroup if there are no processes running
		ioutil.WriteFile(filepath.Join(pids, c.name+"/notify_on_release"), []byte("1"), writeFilePerm),
		ioutil.WriteFile(filepath.Join(pids, c.name+"/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), writeFilePerm),
		// Create memory sub-directory
		os.Mkdir(filepath.Join(memory, c.name), mkdirPerm),
		// Set the maximum memory for the container
		ioutil.WriteFile(filepath.Join(memory, c.name+"memory.limit_in_bytes"), []byte(c.props.memMax), writeFilePerm),
		// Create net_cls sub-directory
		os.Mkdir(filepath.Join(netCls, c.name), mkdirPerm),
		// Set the network id to identify packets from this container
		ioutil.WriteFile(filepath.Join(netCls, c.name+"net_cls.classid"), []byte(string(c.cls)), writeFilePerm),
	)
}
