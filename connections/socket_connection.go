package connections

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type SocketConnection struct {
	Index         int
	LocalAddress  *net.IP
	LocalPort     int
	RemoteAddress *net.IP
	RemotePort    int
	Status        int
	Inode         uintptr
}

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

func (s *SocketConnection) String(socket_type string) string {
	return fmt.Sprintf("%-8s "+
		"Index: %9d "+
		"LAdrs: %15s "+
		"LPort: %5d "+
		"RAdrs: %15s "+
		"RPort: %5d "+
		"Inode: %8d "+
		"Status: %s\n",
		socket_type, s.Index, s.LocalAddress, s.LocalPort, s.RemoteAddress, s.RemotePort, s.Inode, TCPConnectionStatusMap[s.Status])
}

const (
	CH_CTRL_ERR = iota
	CH_CTRL_QUIT
)

type ChannelControl struct {
	MsgType int
	Error   error
}

func parseSocketConnetions(socket_type string, c chan SocketConnection, ctrlChannel chan ChannelControl) {
	fp, err := os.Open(ConnectionSourceMap[socket_type])
	if err != nil {
		ctrlChannel <- ChannelControl{CH_CTRL_ERR, err}
		return
	}
	defer fp.Close()
	s := bufio.NewScanner(fp)
	for n := 0; s.Scan(); n++ {
		if n < 1 {
			continue
		}
		fields := strings.Fields(s.Text())
		t := SocketConnection{}
		if _, err := fmt.Sscanf(fields[0], "%d:", &t.Index); err != nil {
			ctrlChannel <- ChannelControl{CH_CTRL_ERR, err}
			return
		}
		t.LocalAddress, t.LocalPort, err = parseSocketAddress(fields[1])
		if err != nil {
			ctrlChannel <- ChannelControl{CH_CTRL_ERR, err}
			return
		}
		t.RemoteAddress, t.RemotePort, err = parseSocketAddress(fields[2])
		if err != nil {
			ctrlChannel <- ChannelControl{CH_CTRL_ERR, err}
			return
		}
		_, err = fmt.Sscanf(fields[3], "%x", &t.Status)
		if err != nil {
			ctrlChannel <- ChannelControl{CH_CTRL_ERR, err}
			return
		}
		_, err = fmt.Sscanf(fields[9], "%d", &t.Inode)
		if err != nil {
			ctrlChannel <- ChannelControl{CH_CTRL_ERR, err}
			return
		}
		c <- t
	}
	ctrlChannel <- ChannelControl{CH_CTRL_QUIT, nil}
}
