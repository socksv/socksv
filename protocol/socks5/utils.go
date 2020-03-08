package socks5

import (
	"encoding/binary"
	"net"
	"strconv"
)

// ParseAddress format address x.x.x.x:xx to raw address.
// addr contains domain length
func ParseAddress(address string) (a byte, addr []byte, port []byte, err error) {
	var h, p string
	h, p, err = net.SplitHostPort(address)
	if err != nil {
		return
	}
	ip := net.ParseIP(h)
	if ip4 := ip.To4(); ip4 != nil {
		a = ATYPIPv4
		addr = []byte(ip4)
	} else if ip6 := ip.To16(); ip6 != nil {
		a = ATYPIPv6
		addr = []byte(ip6)
	} else {
		a = ATYPDomain
		addr = []byte{byte(len(h))}
		addr = append(addr, []byte(h)...)
	}
	i, _ := strconv.Atoi(p)
	port = make([]byte, 2)
	binary.BigEndian.PutUint16(port, uint16(i))
	return
}
