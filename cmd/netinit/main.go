package main

import (
	"log"
	"flag"

	"github.com/abdulbahajaj/startc/pkg/networking"
)

func main(){
	var pid int
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")
	flag.Parse()

	if pid == 0 {
		log.Panic("PID is not passed to NetInit")
	}

	_, _, err := networking.VethInContainer(pid)
	if err != nil {
		log.Panic(err)
	}
}
