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
)

const (
	mkdirPerm     os.FileMode = 0755
	writeFilePerm os.FileMode = 0700
)

// RunState ... Object describing state of a container
type RunState struct {
	Running  bool
	HasChild bool
}

// Container ... Containerized instance of a linux file system
type Container struct {
	ID      uuid.UUID
	Name    string
	State   RunState
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
	c.State.Running = true
}

// SpawnChild ... Spawn a child process within the fs instance created via run()
func (c *Container) SpawnChild(args []string) {
	log.Printf("Running in new UTS namespace %v as %d\n", args[2:], os.Getpid())

	c.Handler.HandledInvocationGroup(
		syscall.Sethostname([]byte(c.Name)),
		syscall.Chroot("/root/"+c.Props.fsname),
		syscall.Chdir("/"), // set the working directory inside container
		syscall.Mount("proc", "proc", "proc", 0, ""),
	)

	cmd := exec.Command(args[2], args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	c.Handler.HandleErrors(cmd.Run())
	c.Handler.HandleErrors(syscall.Unmount(c.Props.fsname, 0))

	c.State.HasChild = true
}

// CreateCGroup ... Create a CGroup for the spawned process
func (c *Container) CreateCGroup() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	memory := filepath.Join(cgroups, "memory")
	netCls := filepath.Join(cgroups, "net_cls")

	c.Handler.HandledInvocationGroup(
		// Create CGroup sub-directory
		os.Mkdir(filepath.Join(pids, c.Name), mkdirPerm),
		// Set maximum child processes
		ioutil.WriteFile(filepath.Join(pids, c.Name+"/pids.max"), []byte(c.Props.procMax), writeFilePerm),
		// Delete the CGroup if there are no processes running
		ioutil.WriteFile(filepath.Join(pids, c.Name+"/notify_on_release"), []byte("1"), writeFilePerm),
		ioutil.WriteFile(filepath.Join(pids, c.Name+"/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), writeFilePerm),
		// Create memory sub-directory
		os.Mkdir(filepath.Join(memory, c.Name), mkdirPerm),
		// Set the maximum memory for the container
		ioutil.WriteFile(filepath.Join(memory, c.Name+"memory.limit_in_bytes"), []byte(c.Props.memMax), writeFilePerm),
		// Create net_cls sub-directory
		os.Mkdir(filepath.Join(netCls, c.Name), mkdirPerm),
		// Set the network id to identify packets from this container
		ioutil.WriteFile(filepath.Join(netCls, c.Name+"net_cls.classid"), []byte(string(c.Cls)), writeFilePerm),
	)
}

// LoadDiskImage ... Load an ISO image to create a containerized image
func (c *Container) LoadDiskImage(diskImg string) {
	c.Im.CreateIso(diskImg)
	log.Println("Load disk image: [%s]", c.Im.VolumeLabel)
}
