package namespaces

import (
	"os"
	"log"
	"syscall"
	"os/exec"
	"fmt"

	"github.com/moby/moby/pkg/reexec"

	// "github.com/abdulbahajaj/startc/pkg/cgroup"
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

	CPUShare int
	MemoryLimit int

	// Other namespace settings
	Persistent string    // Path to where to save dirs
	MountProc bool       // Whether or not to mount the proc filesystem
	Cmd string           // The command to run in the new namespace
	MountPath string           // The command to run in the new namespace
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

	if ns.MountPath == "" {
		log.Fatal("Provide a valid mount path")
	}

	cmd := reexec.Command("nsInit", "--cmd", ns.Cmd, "--mount", ns.MountPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: getFlags(ns),
		UidMappings: getMapping(0, os.Getuid(), 1),
		GidMappings: getMapping(0, os.Getgid(), 1),
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Created namespace %+v\n", ns)

	pid := fmt.Sprintf("%d", cmd.Process.Pid)
	netinitCmd := exec.Command("bin/pinit", "-pid", pid)
	netinitCmd.Stdin = os.Stdin
	netinitCmd.Stdout = os.Stdout
	netinitCmd.Stderr = os.Stderr

	if err := netinitCmd.Run(); err != nil {
		log.Panic(err)
	}


	cmd.Process.Signal(syscall.SIGIO)

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Println("Exited container")
}

func init() {
	reexec.Register("nsInit", nsInit)
	if reexec.Init() {
		os.Exit(0)
	}
}
