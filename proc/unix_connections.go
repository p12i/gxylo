package proc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var UnixConnectionStatusMap = map[int]string{
	0x00: "CLOSE",
	0x01: "SYN_SENT",
	0x02: "ESTABLISHED",
	0x03: "CLOSING",
}

type UnixConnection struct {
	Number   int
	RefCount int
	Protocol int
	Flags    int
	UnixType int
	State    int
	Inode    uintptr
	Path     string
}

func (t *UnixConnection) String() string {
	return fmt.Sprintf("Type:           unix\n"+
		"Number:         %d\n"+
		"RefCount:       %d\n"+
		"Protocol:       %d\n"+
		"Flags:          %d\n"+
		"UnixSocketType: %d\n"+
		"State:          %s\n"+
		"Inode:          %d\n"+
		"Path:           %s\n", t.Number, t.RefCount, t.Protocol, t.Flags, t.UnixType, UnixConnectionStatusMap[t.State], t.Inode, t.Path)

}

// Num               RefCount Protocol Flags    Type St Inode Path
// 0000000000000000: 00000002 00000000 00010000 0001 01 37071 @/tmp/dbus-uySDV1kH
// 0000000000000000: 00000002 00000000 00010000 0001 01 25105 /run/uuidd/request
// 0000000000000000: 00000002 00000000 00010000 0001 01 25109 /var/run/dbus/system_bus_socket

func (l *ConnectionList) ParseUnixConnections() error {
	fp, err := os.Open(ConnectionSourceMap["unix"])
	if err != nil {
		return err
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	for n := 0; s.Scan(); n++ {
		if n < 1 {
			continue
		}
		fields := strings.Fields(s.Text())
		t := UnixConnection{}
		_, err := fmt.Sscanf(fields[0], "%x:", &t.Number)
		if err != nil {
			return err
		}
		if _, err = fmt.Sscanf(fields[1], "%x", &t.RefCount); err != nil {
			return err
		}
		if _, err = fmt.Sscanf(fields[2], "%x", &t.Protocol); err != nil {
			return err
		}
		if _, err = fmt.Sscanf(fields[3], "%x", &t.Flags); err != nil {
			return err
		}

		if _, err = fmt.Sscanf(fields[4], "%x", &t.UnixType); err != nil {
			return err
		}

		if _, err = fmt.Sscanf(fields[5], "%x", &t.State); err != nil {
			return err
		}

		if _, err = fmt.Sscanf(fields[6], "%d", &t.Inode); err != nil {
			return err
		}
		if len(fields) == 8 {
			t.Path = fields[7]
		} else {
			t.Path = ""
		}
		l.Connections[t.Inode] = &t
	}

	return nil
}
