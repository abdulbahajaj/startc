package main

import (
	"log"
	"flag"
	"net"

	"github.com/abdulbahajaj/startc/pkg/networking"
	// "github.com/vishvananda/netlink"
)

func main(){
	var pid int
	flag.IntVar(&pid, "pid", 0, "pid of a process in the container's network namespace")
	flag.Parse()

	if pid == 0 {
		log.Panic("PID is not passed to NetInit")
	}

	bridge, err := networking.GetDefaultBridge()
	if err != nil {
		log.Panic(err)
	}

	veth1, veth2, err := networking.CreateVeth()

	networking.AttachVethToBridge(veth1, bridge)

	networking.ContainerSetup(pid, net.IPv4(172, 0, 0, 2), net.IPv4Mask(255, 255, 0, 0), veth2)

	// netlink.LinkSetNsPid(veth2, pid)







	// bridge, err := networking.CreateBridge()
	// if err != nil {
		// log.Panic(err)
	// }

	// networking.AttachToContainer(bridge, pid)

	// // _, _, err := networking.VethInContainer(pid)
	// if err != nil {
		// log.Panic(err)
	// }
}
