package main

import (
	"log"
	"flag"

	"github.com/abdulbahajaj/startc/pkg/networking"
	"github.com/abdulbahajaj/startc/pkg/cgroup"
)

func main(){
	var pid int
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")

	flag.Parse()

	if pid == 0 {
		log.Panic("PID is not passed to NetInit")
	}


	br, err := networking.EnvSetup()
	if err != nil {
		log.Panic(err)
	}

	cont, err := networking.ContainerSetup(pid, 1, br)

	if err != nil {
		log.Panic(err)
	}

	log.Println(cont)

	// if (ns.Cgroup){
	cgDesc := cgroup.Desc{
		PID: pid,
		CPUShare: 1010,
		MemoryLimit: 101000000,
	}
	if err := cgroup.Apply(cgDesc); err != nil {
		log.Panic(err)
	}

// }

}
