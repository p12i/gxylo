package proc

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
	"raw":     "/proc/net/raw",
	"raw6":    "/proc/net/raw6",
	"packet":  "/proc/net/packet",
	"netlink": "/proc/net/netlink",
	"unix":    "/proc/net/unix",
}

type Connection interface {
	String() string
}

type ConnectionList struct {
	Connections map[uintptr]Connection
}

func (l *ConnectionList) ParseConnections() error {
	l.Connections = make(map[uintptr]Connection)

	if err := l.ParseTCPConnections(); err != nil {
		return err
	}

	if err := l.ParseUDPConnections(); err != nil {
		return err
	}
	if err := l.ParseUnixConnections(); err != nil {
		return err
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
	ip = net.IPv4(byte_ip[3], byte_ip[2], byte_ip[1], byte_ip[0])
	return &ip, port, nil

}

func (l *ConnectionList) GetConnection(p uintptr) Connection {
	return l.Connections[p]
}
