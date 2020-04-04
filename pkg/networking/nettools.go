package networking

import (
	"net"
	"time"
	"fmt"
	"strconv"
	"math/rand"
	"runtime"
	"errors"
	"os/exec"

	"github.com/vishvananda/netns"
	"github.com/vishvananda/netlink"
)

type Bridge struct {
	Name string
	Link netlink.Link
	IP net.IP
	Mask net.IPMask
}

type Container struct {
	PID int
	IP net.IP
	InVeth netlink.Link
	OutVeth netlink.Link
	Bridge Bridge
}

func makeInterfaceName(label string) string{
	rand.Seed(time.Now().UTC().UnixNano())
	return fmt.Sprintf("sc-%s%s", label, strconv.Itoa(rand.Int())[:6])
}

func CreateBridge(name string, ip net.IP, mask net.IPMask) (Bridge, error) {
	linkAttrs := netlink.LinkAttrs{Name: name}
	netBridge := netlink.Bridge{LinkAttrs: linkAttrs}

	if err := netlink.LinkAdd(&netBridge); err != nil {
		return Bridge{}, err
	}

	address := &netlink.Addr{IPNet: &net.IPNet{IP: ip, Mask: mask}}
	if err := netlink.AddrAdd(&netBridge, address); err != nil {
		return Bridge{}, err
	}

	if err := netlink.LinkSetUp(&netBridge); err != nil {
		return Bridge{}, err
	}

	link, err := netlink.LinkByName(name)
	if err != nil {
		return Bridge{}, err
	}

	return Bridge{
		Name: name,
		Link: link,
		IP: ip,
		Mask: mask,
	}, nil
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

func AddDefaultRoute(ip net.IP, dev netlink.Link) error {
		route := &netlink.Route{
			Scope:     netlink.SCOPE_UNIVERSE,
			LinkIndex: dev.Attrs().Index,
			Gw:        ip,
		}

		return netlink.RouteAdd(route)
}

func CalculateContainerIP(cid int) (net.IP, error) {
	c, d := 0, cid + 1

	for d > 255 {
		d = d - 255
		c += 1
	}

	if c > 255 {
		return nil, errors.New("Exceeded maximum allowed IP range")
	}

	return net.IPv4(172, 0, byte(c), byte(d)), nil
}

func ContainerSetup(pid int, cid int, bridge Bridge) (Container, error) {
	veth1, veth2, err := CreateVeth()
	AttachVethToBridge(veth1, bridge.Link)
	netlink.LinkSetNsPid(veth2, pid)
	containerIP, err := CalculateContainerIP(cid)

	if err != nil {
		return Container{}, err
	}

	if err := ExecInNS(pid, func() error {
		addr := &netlink.Addr{IPNet: &net.IPNet{IP: containerIP, Mask: bridge.Mask}}
		if err := netlink.AddrAdd(veth2, addr); err != nil {
			return err
		}

		if err := netlink.LinkSetUp(veth2); err != nil {
			return err
		}

		return AddDefaultRoute(bridge.IP, veth2)
	}); err != nil {
		return Container{}, err
	}

	return Container{
		PID: pid,
		IP: containerIP,
		InVeth: veth2,
		OutVeth: veth1,
		Bridge: bridge,
	}, nil
}

func GetCIDR(ip net.IP, mask net.IPMask) string {
	binary := strconv.FormatInt(int64(mask[0]), 2)
	binary += strconv.FormatInt(int64(mask[1]), 2)
	binary += strconv.FormatInt(int64(mask[2]), 2)
	binary += strconv.FormatInt(int64(mask[3]), 2)

	count := 0
	for cur := 0; cur < len(binary); cur++ {
		if binary[cur] == '1' {
			count += 1
		}
	}

	return ip.String() + "/" + strconv.Itoa(count)
}

func EnableIPForward() error {
	return exec.Command("sysctl", "net.ipv4.ip_forward=1").Run()
}

// func GetDefaultGateway() err {
// 	routeOut, err := netlink.RouteGet(net.IPv4(1,1,1,1))
// 	if err != nil {
// 		return err
// 	}
// 	if len(routeOut) == 0 {
// 		return errors.New("No route out was found")
// 	}

// 	ifindex := routeOut[0].LinkIndex

// 	links, err :=  netlink.LinkList()
// 	if err != nil {
// 		return err
// 	}

// 	if len(links) < ifindex {
// 		return errors.New("Was not able to retrieve desired link")
// 	}

// 	links[ifindex - 1].Attrs().Name
// }

func Masquerade(bridge Bridge) error {
	return exec.Command("iptables",
			"-t", "nat",
			"-A", "POSTROUTING",
			"-s", GetCIDR(bridge.IP, bridge.Mask),
			"-j", "MASQUERADE").Run()
}

func EnvSetup() (Bridge, error) {
	name := "sc-br0"
	mask := net.IPv4Mask(255, 255, 0, 0)
	bridgeIP := net.IPv4(172, 0, 0, 1)

	if link, err := netlink.LinkByName(name); err == nil {
		return Bridge{
			Name: name,
			Link: link,
			Mask: mask,
			IP: bridgeIP,
		}, err
	}

	bridge, err := CreateBridge(name, bridgeIP, mask)
	if err != nil {
		return Bridge{}, err
	}
	if err := EnableIPForward(); err != nil {
		return Bridge{}, err
	}
	if err := Masquerade(bridge); err != nil {
		return Bridge{}, err
	}

	return bridge, nil
}
