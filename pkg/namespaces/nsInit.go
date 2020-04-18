package namespaces

import (
	"log"
	"os/exec"
	"os/signal"
	"os"
	"syscall"
	"path/filepath"
	"flag"
)

func setSysfs(path string) error {
	os.MkdirAll("/sys", 0755)
	source := ""
	target := filepath.Join(path, "/sys")
	fstype := "sysfs"
	flags := 0
	data := ""
	return syscall.Mount(source, target, fstype, uintptr(flags), data)
}

func setProc(path string) error {
	os.MkdirAll("/proc", 0755)
	source := ""
	target := filepath.Join(path, "/proc")
	fstype := "proc"
	flags := 0
	data := ""
	return syscall.Mount(source, target, fstype, uintptr(flags), data)
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

func run(cmd string) error {
	log.Println("Running container")
	process := exec.Command("/bin/sh")
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	return process.Run()
}

func WaitUntilSetup() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGIO)
	<-sigs
}

/*
* Entrypoint to the namespace
*/

func nsInit() {
	var cmd string
	var mountPath string

	flag.StringVar(&cmd, "cmd", "/bin/sh", "Enter the command that you want to run inside the container")
	flag.StringVar(&mountPath, "mount", "", "Enter the command that you want to run inside the container")
	flag.Parse()

	if mountPath == "" {
		log.Panic("You should provide a mount path")
	}


	log.Println("Initializing namespaces")
	if err := setProc(mountPath); err != nil {
		log.Panic(err)
	}

	if err := setSysfs(mountPath); err != nil {
		log.Panic(err)
	}

	if err := setRoot(mountPath); err != nil {
		log.Panic(err)
	}
	if err := syscall.Sethostname([]byte("startc")); err != nil {
		log.Panic(err)
	}


	WaitUntilSetup()

	if err := run(cmd); err != nil {
		log.Panic(err)
	}
}
