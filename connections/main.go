package connections

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"
)

var ConnectionSourceMap = map[string]string{
	"tcp":     "/proc/net/tcp",
	"tcp6":    "/proc/net/tcp6",
	"udp":     "/proc/net/udp",
	"udp6":    "/proc/net/udp6",
	"unix":    "/proc/net/unix",
	"raw":     "/proc/net/raw",
	"raw6":    "/proc/net/raw6",
	"packet":  "/proc/net/packet",  // TODO
	"netlink": "/proc/net/netlink", // TODO
}

type Connection interface {
	String() string
}



type ConnectionList struct {
	Connections map[uintptr]Connection
}

func (l *ConnectionList) ParseConnections() error {
	l.Connections = make(map[uintptr]Connection)
	var functions = []func() error{
		l.ParseTCPConnections,
		l.ParseTCP6Connections,
		l.ParseUDPConnections,
		l.ParseUDP6Connections,
		l.ParseUnixConnections,
		l.ParseRawConnections,
		l.ParseRaw6Connections,
	}

	for _, f := range functions {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func (l *ConnectionList) String() string {
	var str strings.Builder
	for _, elem := range l.Connections {
		str.WriteString(elem.String())
	}
	return str.String()
}

func parseSocketAddress(s string) (*net.IP, int, error) {
	endpoint := strings.Split(s, ":")
	ip := net.IP{}
	port := 0
	if len(endpoint) != 2 {
		return &ip, 0, errors.New(fmt.Sprintf("Invalid number of slices for "+s+". Expected 2 got %d", len(endpoint)))
	}

	byte_ip, err := hex.DecodeString(endpoint[0])
	if err != nil {
		return &ip, port, err
	}
	_, err = fmt.Sscanf(endpoint[1], "%x", &port)
	if err != nil {
		return &ip, port, err
	}
	if len(byte_ip) == 4 {
		ip = net.IPv4(byte_ip[3], byte_ip[2], byte_ip[1], byte_ip[0])
	} else if len(byte_ip) == 16 {
		ip = net.IP{
			byte_ip[15], byte_ip[14], byte_ip[13], byte_ip[12],
			byte_ip[11], byte_ip[10], byte_ip[9], byte_ip[8],
			byte_ip[7], byte_ip[6], byte_ip[5], byte_ip[4],
			byte_ip[3], byte_ip[2], byte_ip[1], byte_ip[0]}
	}

	return &ip, port, nil

}

func (l *ConnectionList) GetConnection(p uintptr) Connection {
	return l.Connections[p]
}
