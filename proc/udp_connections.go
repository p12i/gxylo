package proc

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type UDPConnection struct {
	Index         int
	LocalAddress  *net.IP
	LocalPort     int
	RemoteAddress *net.IP
	RemotePort    int
	TCPStatus     int
	Inode         uintptr
}

func (t *UDPConnection) String() string {
	return fmt.Sprintf("Type:          udp\n"+
		"Index:         %d\n"+
		"LocalAddress:  %s\n"+
		"LocalPort:     %d\n"+
		"RemoteAddress: %s\n"+
		"RemotePort:    %d\n"+
		"TCPStatus:     %s\n"+
		"Inode:         %d\n", t.Index, t.LocalAddress, t.LocalPort, t.RemoteAddress, t.RemotePort, TCPConnectionStatusMap[t.TCPStatus], t.Inode)

}

//   sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
//   76: 0103000A:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 41260 2 0000000000000000 0
//      76: 3500007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000   101        0 26015 2 0000000000000000 0
//
func (l *ConnectionList) ParseUDPConnections() error {
	fp, err := os.Open(ConnectionSourceMap["udp"])
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
		t := UDPConnection{}
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
