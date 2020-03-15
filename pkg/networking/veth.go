package networking

import (
	"runtime"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)


func CreateVethPair() (netlink.Link, netlink.Link, error) {

	containerInterfaceName := makeInterfaceName()
	hostInterfaceName :=  makeInterfaceName()

	attrs := netlink.NewLinkAttrs()
	attrs.Name = containerInterfaceName


	veth := &netlink.Veth{
		LinkAttrs: attrs,
		PeerName: hostInterfaceName,
	}


	if err := netlink.LinkAdd(veth); err != nil {
		return nil, nil, err
	}

	if err := netlink.LinkSetUp(veth); err != nil {
		return nil, nil, err
	}

	peer, err := netlink.LinkByName(hostInterfaceName)
	if err != nil {
		return nil, nil, err
	}

	if err := netlink.LinkSetUp(veth); err != nil {
		return nil, nil, err
	}

	return veth, peer, nil
}

func VethInContainer(pid int) (netlink.Link, netlink.Link, error) {
	veth, peer, err := CreateVethPair()
	if err != nil {
		return nil, nil, err
	}


	if err := netlink.LinkSetNsPid(peer, pid); err != nil {
		// TODO delete created interfaces in case of failure
		return nil, nil, err
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	currentNs, err := netns.Get()
	if err != nil {
		return nil, nil, err
	}
	defer currentNs.Close()

	nnsHandler, err := netns.GetFromPid(pid)
	if err != nil {
		return nil, nil, err
	}
	defer nnsHandler.Close()

	if err := netns.Set(nnsHandler); err != nil {
		return nil, nil, err
	}


	if err := netlink.LinkSetUp(peer); err != nil {
		return nil, nil, err
	}

	return veth, peer, err
}

