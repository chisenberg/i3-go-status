package netblock

import (
	"net"

	"github.com/chisenberg/i3-go-status/block"
)

// Provider shows the first non-loopback IPv4 unicast address (skips link-local), interface order from the OS.
type Provider struct{}

// New returns a BlockInterface for the local network IP block.
func New() *Provider {
	return &Provider{}
}

// GetBlock implements block.BlockInterface.
func (*Provider) GetBlock() *block.Block {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ip := addrToIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				if ip4.IsLinkLocalUnicast() {
					continue
				}
				return &block.Block{
					Name:     "network",
					FullText: ip4.String(),
				}
			}
		}
	}
	return nil
}

// ClickBlock implements block.BlockInterface.
func (*Provider) ClickBlock(block.ClickEvent) {}

func addrToIP(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	default:
		return nil
	}
}
