package namespaces

import (
	"os"
	"log"
	"syscall"

	"github.com/moby/moby/pkg/reexec"
)

type Desc struct {
	// What namespaces to create
	Mount bool
	Uts bool
	Ipc bool
	Network bool
	Pid bool
	Cgroup bool
	User bool

	// Other namespace settings
	Persistent string    // Path to where to save dirs
	MountProc bool       // Whether or not to mount the proc filesystem
	Cmd string           // The command to run in the new namespace
}

func orIf(cond bool, flag uintptr, condFlag uintptr) uintptr{
	if cond {
		flag |= condFlag
	}
	return flag
}

func getFlags(ns Desc) uintptr {
	var flags uintptr = 0x1
	flags = orIf(ns.Ipc, flags, syscall.CLONE_NEWIPC)
	flags = orIf(ns.Network, flags, syscall.CLONE_NEWNET)
	flags = orIf(ns.Mount, flags, syscall.CLONE_NEWNS)
	flags = orIf(ns.Pid, flags, syscall.CLONE_NEWPID)
	flags = orIf(ns.User, flags, syscall.CLONE_NEWUSER)
	flags = orIf(ns.Uts, flags, syscall.CLONE_NEWUTS)
	return flags
}

func getMapping(containerID int, hostID int, size int)[]syscall.SysProcIDMap{
	return []syscall.SysProcIDMap{
		{
			ContainerID: containerID,
			HostID: hostID,
			Size: size,
		},
	}
}

func Create(ns Desc){

	cmd := reexec.Command("nsInit")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: getFlags(ns),
		UidMappings: getMapping(0, os.Getuid(), 1),
		GidMappings: getMapping(0, os.Getgid(), 1),
	}

	log.Printf("Creating a new namespace %+v\n", ns)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Exited container")
}

func init() {
	reexec.Register("nsInit", nsInit)
	if reexec.Init() {
		os.Exit(0)
	}
}
