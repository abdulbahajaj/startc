package networking

import (
	"net"
	"time"
	"fmt"
	"strconv"
	"math/rand"
	"runtime"

	"github.com/vishvananda/netns"
	"github.com/vishvananda/netlink"
)

func makeInterfaceName(label string) string{
	rand.Seed(time.Now().UTC().UnixNano())
	return fmt.Sprintf("sc-%s%s", label, strconv.Itoa(rand.Int())[:6])
}

func GetDefaultBridge() (netlink.Link, error) {
	name := "sc-br0"
	if bridge, err := netlink.LinkByName(name); err == nil {
		return bridge, err
	}

	ip := net.IPv4(172, 0, 0, 1)
	mask := net.IPv4Mask(255, 255, 0, 0)
	return CreateBridge(name, ip, mask)
}

func CreateBridge(name string, ip net.IP, mask net.IPMask) (netlink.Link, error) {
	linkAttrs := netlink.LinkAttrs{Name: name}
	bridge := netlink.Bridge{LinkAttrs: linkAttrs}

	if err := netlink.LinkAdd(&bridge); err != nil {
		return nil, err
	}

	address := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: mask}}
	if err := netlink.AddrAdd(&bridge, address); err != nil {
		return nil, err
	}

	if err := netlink.LinkSetUp(&bridge); err != nil {
		return nil, err
	}

	return netlink.LinkByName(name)
}

func CreateVeth() (netlink.Link, netlink.Link, error){
	invName := makeInterfaceName("eth")
	outvName := makeInterfaceName("eth")

	attrs := netlink.NewLinkAttrs()
	attrs.Name = outvName
	outv := &netlink.Veth{
		LinkAttrs: attrs,
		PeerName: invName,
	}


	if err := netlink.LinkAdd(outv); err != nil {
		return nil, nil, err
	}

	if err := netlink.LinkSetUp(outv); err != nil {
		return nil, nil, err
	}


	inv, err := netlink.LinkByName(invName)
	if err != nil {
		return nil, nil, err
	}

	if err := netlink.LinkSetUp(inv); err != nil {
		return nil, nil, err
	}


	return inv, outv, nil
}

func AttachVethToBridge(veth  netlink.Link, bridge netlink.Link) error {
	return netlink.LinkSetMaster(veth, bridge)
}

func MoveVethToContainer(veth netlink.Link, pid int) error {
	if err := netlink.LinkSetNsPid(veth, pid); err != nil {
		return err
	}

	// runtime.LockOSThread()
	// defer runtime.UnlockOSThread()

	// currentNs, err := netns.Get()
	// if err != nil {
		// return err
	// }
	// defer currentNs.Close()
	// defer netns.Set(currentNs)

	// nnsHandler, err := netns.GetFromPid(pid)
	// if err != nil {
		// return err
	// }
	// defer nnsHandler.Close()

	// if err := netns.Set(nnsHandler); err != nil {
		// return err
	// }
	//
	// if err := netlink.LinkSetUp(veth); err != nil {
		// return err
	// }

	return nil
}

func ExecInNS(pid int, fn func() error ) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := netns.Get()
	if err != nil {
		return err
	}
	defer currentNs.Close()
	defer netns.Set(currentNs)

	nnsHandler, err := netns.GetFromPid(pid)
	if err != nil {
		return err
	}
	defer nnsHandler.Close()


	if err := netns.Set(nnsHandler); err != nil {
		return err
	}

	return fn()
}

func ContainerSetup(pid int, ip net.IP, mask net.IPMask, veth netlink.Link) error {
	netlink.LinkSetNsPid(veth, pid)

	if err := ExecInNS(pid, func() error {
		addr := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: mask}}
		if err := netlink.AddrAdd(veth, addr); err != nil {
			return err
		}

		if err := netlink.LinkSetUp(veth); err != nil {
			return err
		}



		return nil
	}); err != nil {
		return err
	}

	return nil
}



// func SetRoute(ip net.IP, mask net.IPMask, gateway net.IP, dev netlink.Link) error {

// }

// func AddAddr(ip net.IP, mask net.IPMask, dev netlink.Link) error {

// }

// func EnablePortForwarding() error {

// }

// func Masquerade(source net.IP, output netlink.Link){

// }
