package proc

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var TCPConnectionStatusMap = map[int]string{
	0x01: "ESTABLISHED",
	0x02: "SYN_SENT",
	0x03: "SYN_RECV",
	0x04: "FIN_WAIT1",
	0x05: "FIN_WAIT2",
	0x06: "TIME_WAIT",
	0x07: "CLOSE",
	0x08: "CLOSE_WAIT",
	0x09: "LAST_ACK",
	0x0A: "LISTEN",
	0x0B: "CLOSING",
}

type TCPConnection struct {
	Index         int
	LocalAddress  *net.IP
	LocalPort     int
	RemoteAddress *net.IP
	RemotePort    int
	TCPStatus     int
	Inode         uintptr
}

func (t *TCPConnection) String() string {
	return fmt.Sprintf("Type:          tcp\n"+
		"Index:         %d\n"+
		"LocalAddress:  %s\n"+
		"LocalPort:     %d\n"+
		"RemoteAddress: %s\n"+
		"RemotePort:    %d\n"+
		"TCPStatus:     %s\n"+
		"Inode:         %d\n", t.Index, t.LocalAddress, t.LocalPort, t.RemoteAddress, t.RemotePort, TCPConnectionStatusMap[t.TCPStatus], t.Inode)

}

//  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
//   0: 0103000A:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 40116 1 0000000000000000 100 0 0 10 0
func (l *ConnectionList) ParseTCPConnections() error {
	fp, err := os.Open(ConnectionSourceMap["tcp"])
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
		t := TCPConnection{}
		_, err := fmt.Sscanf(fields[0], "%d:", &t.Index)

		if err != nil {
			return err
		}
		t.LocalAddress, t.LocalPort, err = parseSocketAddress(fields[1])
		if err != nil {
			return err
		}
		t.RemoteAddress, t.RemotePort, err = parseSocketAddress(fields[2])
		if err != nil {
			return err
		}
		_, err = fmt.Sscanf(fields[3], "%x", &t.TCPStatus)
		if err != nil {
			return err
		}
		_, err = fmt.Sscanf(fields[9], "%d", &t.Inode)
		if err != nil {
			return err
		}
		l.Connections[t.Inode] = &t
	}

	return nil
}
