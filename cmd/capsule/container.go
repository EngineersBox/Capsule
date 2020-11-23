package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/google/uuid"
)

// Container ... Containerized instance of a linux file system
type Container struct {
	id   uuid.UUID
	name string
}

// Run ... Fork-execute a fs instance
func (c *Container) Run(args []string) {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args[2:]...)...)

	c.CreateCGroup()

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	handleErrors(cmd.Run())
}

// SpawnChild ... Spawn a child process within the fs instance created via run()
func (c *Container) SpawnChild(args []string) {
	handleErrors(syscall.Mount(props.fsname, props.fsname, "", syscall.MS_BIND, ""))
	handleErrors(os.MkdirAll(props.fsname+"/oldrootfs", 0700))
	handleErrors(syscall.PivotRoot(props.fsname, props.fsname+"/oldrootfs"))
	handleErrors(os.Chdir("/"))

	cmd := exec.Command(args[2], args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	handleErrors(cmd.Run())
	handleErrors(syscall.Unmount(props.fsname, 0))
}

// CreateCGroup ... Create a CGroup for the spawned process
func (c *Container) CreateCGroup() {
	cgroups := "/sys/fs/cgroups"
	pids := filepath.Join(cgroups, "pids")

	handleErrors(os.Mkdir(filepath.Join(pids, c.name), 0755))
	handleErrors(ioutil.WriteFile(filepath.Join(pids, "pids/"+c.name+".max"), []byte(props.procMax), 0700))
	handleErrors(ioutil.WriteFile(filepath.Join(pids, "pids/notify_on_release"), []byte("1"), 0700))
	handleErrors(ioutil.WriteFile(filepath.Join(pids, "pids/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}
