package namespaces

import (
	"fmt"
	"os"
	"log"
	"os/exec"
	"syscall"
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
	Persistent string    // Path to where to save dirs.
	Fork bool            // Whether or not to fork the proccess before exec.
	MountProc bool       // Whether or not to mount the proc filesystem
	Cmd string           // The command to run in the new namespace
}

func Create(ns Desc){
	log.Printf("Creating a new namespace %+v\n", ns)

	// cmd := make([]string, 0)
	// out, err := exec.Command("unshare", cmd...).Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

	// log.Println(unshareFlags)

	// err1 := syscall.Unshare(unshareFlags)

	// if err1 != nil {
	// 	log.Fatal(err1)
	// }

	// log.Println(out)
}

	// unshareFlags := syscall.CLONE_FILES
	// unshareFlags |=	syscall.CLONE_FS
	// unshareFlags |=	syscall.CLONE_NEWCGROUP
	// unshareFlags |=	syscall.CLONE_NEWIPC
	// unshareFlags |=	syscall.CLONE_NEWNET
	// unshareFlags |=	syscall.CLONE_NEWNS
	// unshareFlags |=	syscall.CLONE_NEWPID
	// unshareFlags |=	syscall.CLONE_NEWUSER
	// unshareFlags |=	syscall.CLONE_NEWUTS
