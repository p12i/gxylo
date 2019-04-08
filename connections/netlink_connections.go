package connections

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type NetlinkConnection struct {
	Sk     int
	Eth    int
	Pid    uintptr
	Groups int64
	Rmem   int
	Wmem   int
	Dump   int
	Locks  int
	Drops  int
	Inode  uintptr
}

func (t NetlinkConnection) String() string {
	return fmt.Sprintf(
		"%-8s "+
			"Eth: %3d "+
			"Pid: %12d "+
			"Groups: %5x "+
			"Rmem: %3d "+
			"Wmem: %3d "+
			"Dump: %3d "+
			"Locks: %3d "+
			"Drops: %3d "+
			"Inode: %8d\n",
		"Netlink", t.Eth, t.Pid, t.Groups, t.Rmem, t.Wmem, t.Dump, t.Locks, t.Drops, t.Inode)

}

// Num               RefCount Protocol Flags    Type St Inode Path
// 0000000000000000: 00000002 00000000 00010000 0001 01 37071 @/tmp/dbus-uySDV1kH
// 0000000000000000: 00000002 00000000 00010000 0001 01 25105 /run/uuidd/request
// 0000000000000000: 00000002 00000000 00010000 0001 01 25109 /var/run/dbus/system_bus_socket

func (l *ConnectionList) ParseNetlinkConnections() error {
	fp, err := os.Open(ConnectionSourceMap["netlink"])
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
		t := NetlinkConnection{}

		for _, elem := range []tParsingStruct{
			tParsingStruct{0, "%x", &t.Sk},
			tParsingStruct{1, "%d", &t.Eth},
			tParsingStruct{2, "%d", &t.Pid},
			tParsingStruct{3, "%x", &t.Groups},
			tParsingStruct{4, "%d", &t.Rmem},
			tParsingStruct{5, "%d", &t.Wmem},
			tParsingStruct{6, "%d", &t.Dump},
			tParsingStruct{7, "%d", &t.Locks},
			tParsingStruct{8, "%d", &t.Drops},
			tParsingStruct{9, "%d", &t.Inode},
		} {
			if _, err := fmt.Sscanf(fields[elem.Field], elem.Format, elem.Pointer); err != nil {
				fmt.Println("Error ", err)
				return err
			}
		}
		l.Connections[t.Inode] = &t
	}

	return nil
}
