package cgroup

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"strconv"
	"fmt"
	"log"
)

type Desc struct {
	PID int
	CPUShare int
	MemoryLimit int
	CGroupPath string
}

type CGPaths struct {
	CPU string
	Mem string
}

func getCGPaths(cg Desc) (CGPaths, error) {
	if cg.CGroupPath == "" {
		cg.CGroupPath = "/sys/fs/cgroup/"
	}

	pidStr := fmt.Sprintf("%d",cg.PID)
	mem := filepath.Join(cg.CGroupPath, "memory", pidStr)
	cpu := filepath.Join(cg.CGroupPath, "cpu", pidStr)


	log.Println("Creating cgroup if they don't exist")
	_ = os.Mkdir(cpu, 0655)
	_ = os.Mkdir(mem, 0655)
	// if err := os.Mkdir(mem, 0655); err != nil && err != os.ErrExist {
	// 	return CGPaths{}, err
	// }
	// if err := os.Mkdir(cpu, 0655); err != nil && err != os.ErrExist {
	// 	return CGPaths{}, err
	// }
	log.Println("Created cgroup if they don't exist")

	return CGPaths{
		CPU: cpu,
		Mem: mem,
	}, nil
}

func NumbToString(num int) []byte{
	return []byte(strconv.Itoa(num))
}

func Apply(cg Desc) error {
	log.Println("Applying cgroup desc")
	cgPath, err := getCGPaths(cg)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(
		filepath.Join(cgPath.CPU, "cpu.shares"),
		NumbToString(cg.CPUShare),
		0644);
	err != nil {
		return err
	}

	if err := ioutil.WriteFile(
		filepath.Join(cgPath.Mem, "memory.limit_in_bytes"),
		NumbToString(cg.MemoryLimit),
		0644);
	err != nil {
		return err
	}

	return AddProc(cg)
}

func AddProc(cg Desc) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from ", r)
		}
	}()
	cgPath, err := getCGPaths(cg)
	if err != nil {
		return err
	}
	procsFiles := [2]string{cgPath.CPU, cgPath.Mem}
	strPid := NumbToString(cg.PID)

	for _, cgPath := range procsFiles{
		cgProcPath := filepath.Join(cgPath, "cgroup.procs")
		ioutil.WriteFile(cgProcPath, strPid, 0644)
	}

	log.Println("Added process to CGroup")
	return nil
}

func Test(){}

func Delete(cg Desc){

}
