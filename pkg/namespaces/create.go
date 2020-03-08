package namespaces

import (
	"os"
	"log"
	"os/exec"
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
	// flags |= syscall.CLONE_FILES
	// flags |= syscall.CLONE_FS
	// flags |= syscall.CLONE_NEWCGROUP

	flags = orIf(ns.Ipc, flags, syscall.CLONE_NEWIPC)
	flags = orIf(ns.Network, flags, syscall.CLONE_NEWNET)
	flags = orIf(ns.Mount, flags, syscall.CLONE_NEWNS)
	flags = orIf(ns.Pid, flags, syscall.CLONE_NEWPID)
	flags = orIf(ns.User, flags, syscall.CLONE_NEWUSER)
	flags = orIf(ns.Uts, flags, syscall.CLONE_NEWUTS)
	return flags
}

func getUserMappings() []syscall.SysProcIDMap{
	return []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID: os.Getuid(),
			Size: 1,
		},
	}
}

func getGroupMappings() []syscall.SysProcIDMap{
	return []syscall.SysProcIDMap{
		{
			ContainerID: 0,
			HostID: os.Getgid(),
			Size: 1,
		},
	}
}

func Create(ns Desc){

	cmd := reexec.Command("nsInit")
	// cmd := exec.Command(ns.Cmd)
	
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: getFlags(ns),
		UidMappings: getUserMappings(),
		GidMappings: getGroupMappings(),
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

func nsInit() {
	log.Println("Initializing namespaces")
	run()
}

func run(){
	log.Println("Starting to run container")
	cmd := exec.Command("/bin/bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}
