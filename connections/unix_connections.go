package connections

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var UnixConnectionStatusMap = map[int]string{
	0x00: "SS_FREE",
	0x01: "SS_UNCONNECTED",
	0x02: "SS_CONNECTING",
	0x03: "SS_CONNECTED",
	0x04: "SS_DISCONNECTING",
}

var UnixSocketType = map[int]string{
	0x01: "STREAM",
	0x02: "DGRAM",
	0x03: "RAW",
	0x04: "RDM",
	0x05: "SEQPACKET",
	0x06: "DCCP",
	0x0A: "PACKET",
}

var UnixSocketFlags = map[int]string{
	1 << 16: "ACC", // SO_ACCEPTCON
	1 << 17: "W",   // SO_WAITDATA
	1 << 18: "N",   // SO_NOSPACE
}

type UnixConnection struct {
	RefCount int
	Protocol int
	Flags    int
	UnixType int
	State    int
	Inode    uint64
	Path     string
}

type tParsingStruct struct {
	Field   int64
	Format  string
	Pointer interface{}
}

func parseFlags(flag int) string {
	var s []string

	for k, v := range UnixSocketFlags {
		if k&flag > 0 {
			s = append(s, v)
		}
	}
	return strings.Join(s, ", ")
}

func (t *UnixConnection) String() string {
	return fmt.Sprintf(
		"%-8s "+
			"RefCount: %6d "+
			"Protocol: %2d "+
			"Flags: %-10s "+
			"USType: %-9s "+
			"State: %-14s "+
			"Inode: %8d "+
			"Path: %s\n", "UNIX", t.RefCount, t.Protocol, parseFlags(t.Flags), UnixSocketType[t.UnixType], UnixConnectionStatusMap[t.State], t.Inode, t.Path)

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

		var null []byte
		for _, elem := range []tParsingStruct{
			tParsingStruct{0, "%x:", &null},
			tParsingStruct{1, "%x", &t.RefCount},
			tParsingStruct{2, "%x", &t.Protocol},
			tParsingStruct{3, "%x", &t.Flags},
			tParsingStruct{4, "%x", &t.UnixType},
			tParsingStruct{5, "%x", &t.State},
			tParsingStruct{6, "%d", &t.Inode},
		} {
			if _, err := fmt.Sscanf(fields[elem.Field], elem.Format, elem.Pointer); err != nil {
				return err
			}
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
