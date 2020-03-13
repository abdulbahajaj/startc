package namespaces

import (
	"log"
	"os/exec"
	"os"
	"syscall"
	"net"

	"path/filepath"
)

func run(){
	log.Println("Running container")
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func setRoot(path string) error {
	putold := filepath.Join(path, "/.pivot_root")
	if err := syscall.Mount(path, path, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err;
	}
	if err := os.MkdirAll(putold, 0700); err != nil {
		return err;
	}
	if err := syscall.PivotRoot(path, putold); err != nil {
		return err;
	}
	if err := os.Chdir("/"); err != nil {
		return err;
	}
	putold = "/.pivot_root"
	if err := syscall.Unmount(putold, syscall.MNT_DETACH); err != nil {
		return err;
	}
	if err := os.RemoveAll(putold); err != nil {
		return err;
	}
	return nil;
}

func setProc(path string) error {
	os.MkdirAll("/proc", 0755)
	source := "proc"
	target := filepath.Join(path, "/proc")
	fstype := "proc"
	flags := 0
	data := ""
	return syscall.Mount(source, target, fstype, uintptr(flags), data)
}

func netInit() error{
	ifs, err := net.Interfaces()

	if err != nil {
		log.Panic(err)
	}

	for _, val := range ifs {
		addrs, err := val.Addrs()
		if err != nil {
			log.Println(err)
			continue
		}

		for _, a := range addrs {
			log.Println(a)
		}

		log.Println(val)
	}

	return nil
}

/*
* Entrypoint to the namespace
*/
func nsInit() {
	log.Println("Initializing namespaces")
	root := "/home/ubuntu/projects/mount-points/newroot"
	if err := setProc(root); err != nil {
		log.Panic(err)
	}

	if err := setRoot(root); err != nil {
		log.Panic(err)
	}

	if err := netInit(); err != nil {
		log.Panic(err)
	}

	run()
}
