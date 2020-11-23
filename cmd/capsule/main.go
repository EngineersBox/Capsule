package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

var (
	cid      = "container"
	fsname   = "rootfs"
	proc_max = "20"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("wat should I do")
	}
}

// fun ... Fork-execute a fs instance
func run() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cgroup()

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

// child ... Spawn a child process within the fs instance created via run()
func child() {
	handleErrors(syscall.Mount(fsname, fsname, "", syscall.MS_BIND, ""))
	handleErrors(os.MkdirAll(fsname+"/oldrootfs", 0700))
	handleErrors(syscall.PivotRoot(fsname, fsname+"/oldrootfs"))
	handleErrors(os.Chdir("/"))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	handleErrors(cmd.Run())
	handleErrors(syscall.Unmount(fsname, 0))
}

// cgroup ... Create a CGroup for the spawned process
func cgroup() {
	cgroups := "/sys/fs/cgroups"
	pids := filepath.Join(cgroups, "pids")

	handleErrors(os.Mkdir(filepath.Join(pids, cid), 0755))
	handleErrors(ioutil.WriteFile(filepath.Join(pids, "pids/"+cid+".max"), []byte(proc_max), 0700))
	handleErrors(ioutil.WriteFile(filepath.Join(pids, "pids/notify_on_release"), []byte("1"), 0700))
	handleErrors(ioutil.WriteFile(filepath.Join(pids, "pids/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}
