package connections

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
)

type PacketConnection struct {
	Sk     int
	RefCnt int
	Type   int
	Proto  int
	Iface  int
	R      int
	Rmem   int
	User   int
	Inode  uintptr
}

//https://godoc.org/golang.org/x/sys/unix
var PacketProtocolType = map[int]string{
	0x88f7: "ETH_P_1588",
	0x88a8: "ETH_P_8021AD",
	0x88e7: "ETH_P_8021AH",
	0x8100: "ETH_P_8021Q",
	0x8917: "ETH_P_80221",
	0x4:    "ETH_P_802_2",
	0x1:    "ETH_P_802_3",
	0x600:  "ETH_P_802_3_MIN",
	0x88b5: "ETH_P_802_EX1",
	0x80f3: "ETH_P_AARP",
	0xfbfb: "ETH_P_AF_IUCV",
	0x3:    "ETH_P_ALL",
	0x88a2: "ETH_P_AOE",
	0x1a:   "ETH_P_ARCNET",
	0x806:  "ETH_P_ARP",
	0x809b: "ETH_P_ATALK",
	0x8884: "ETH_P_ATMFATE",
	0x884c: "ETH_P_ATMMPOA",
	0x2:    "ETH_P_AX25",
	0x4305: "ETH_P_BATMAN",
	0x8ff:  "ETH_P_BPQ",
	0xf7:   "ETH_P_CAIF",
	0xc:    "ETH_P_CAN",
	0xd:    "ETH_P_CANFD",
	0x16:   "ETH_P_CONTROL",
	0x6006: "ETH_P_CUST",
	0x6:    "ETH_P_DDCMP",
	0x6000: "ETH_P_DEC",
	0x6005: "ETH_P_DIAG",
	0x6001: "ETH_P_DNA_DL",
	0x6002: "ETH_P_DNA_RC",
	0x6003: "ETH_P_DNA_RT",
	0x1b:   "ETH_P_DSA",
	0x18:   "ETH_P_ECONET",
	0xdada: "ETH_P_EDSA",
	0x88be: "ETH_P_ERSPAN",
	0x22eb: "ETH_P_ERSPAN2",
	0x8906: "ETH_P_FCOE",
	0x8914: "ETH_P_FIP",
	0x19:   "ETH_P_HDLC",
	0x892f: "ETH_P_HSR",
	0x8915: "ETH_P_IBOE",
	0xf6:   "ETH_P_IEEE802154",
	0xa00:  "ETH_P_IEEEPUP",
	0xa01:  "ETH_P_IEEEPUPAT",
	0xed3e: "ETH_P_IFE",
	0x800:  "ETH_P_IP",
	0x86dd: "ETH_P_IPV6",
	0x8137: "ETH_P_IPX",
	0x17:   "ETH_P_IRDA",
	0x6004: "ETH_P_LAT",
	0x886c: "ETH_P_LINK_CTL",
	0x9:    "ETH_P_LOCALTALK",
	0x60:   "ETH_P_LOOP",
	0x9000: "ETH_P_LOOPBACK",
	0x88e5: "ETH_P_MACSEC",
	0xf9:   "ETH_P_MAP",
	0x15:   "ETH_P_MOBITEX",
	0x8848: "ETH_P_MPLS_MC",
	0x8847: "ETH_P_MPLS_UC",
	0x88f5: "ETH_P_MVRP",
	0x88f8: "ETH_P_NCSI",
	0x894f: "ETH_P_NSH",
	0x888e: "ETH_P_PAE",
	0x8808: "ETH_P_PAUSE",
	0xf5:   "ETH_P_PHONET",
	0x10:   "ETH_P_PPPTALK",
	0x8863: "ETH_P_PPP_DISC",
	0x8:    "ETH_P_PPP_MP",
	0x8864: "ETH_P_PPP_SES",
	0x88c7: "ETH_P_PREAUTH",
	0x88fb: "ETH_P_PRP",
	0x200:  "ETH_P_PUP",
	0x201:  "ETH_P_PUPAT",
	0x9100: "ETH_P_QINQ1",
	0x9200: "ETH_P_QINQ2",
	0x9300: "ETH_P_QINQ3",
	0x8035: "ETH_P_RARP",
	0x6007: "ETH_P_SCA",
	0x8809: "ETH_P_SLOW",
	0x5:    "ETH_P_SNAP",
	0x890d: "ETH_P_TDLS",
	0x6558: "ETH_P_TEB",
	0x88ca: "ETH_P_TIPC",
	0x1c:   "ETH_P_TRAILER",
	0x11:   "ETH_P_TR_802_2",
	0x22f0: "ETH_P_TSN",
	0x7:    "ETH_P_WAN_PPP",
	0x883e: "ETH_P_WCCP",
	0x805:  "ETH_P_X25",
	0xf8:   "ETH_P_XDSA",
}

func (t PacketConnection) String() string {
	var pInterface string
	var pUser string
	var proto string
	var ok bool

	if iface, err := net.InterfaceByIndex(t.Iface); err != nil {
		pInterface = strconv.Itoa(t.Iface)
	} else {
		pInterface = iface.Name
	}

	if lUser, err := user.LookupId(strconv.Itoa(t.User)); err != nil {
		pUser = strconv.Itoa(t.User)
	} else {
		pUser = lUser.Username
	}

	if proto, ok = PacketProtocolType[t.Proto]; !ok {
		proto = fmt.Sprintf("%x", t.Proto)
	}

	return fmt.Sprintf(
		"%-8s "+
			"RefCnt: %3d "+
			"Type: %-9s "+
			"Proto: %-12s "+
			"Iface: %10s "+
			"R: %3d "+
			"Rmem: %3d "+
			"User: %12s "+
			"Inode: %8d\n",
		"Packet", t.RefCnt, UnixSocketType[t.Type], proto, pInterface, t.R, t.Rmem, pUser, t.Inode)

}

func (l *ConnectionList) ParsePacketConnections() error {
	fp, err := os.Open(ConnectionSourceMap["packet"])
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
		t := PacketConnection{}

		for _, elem := range []tParsingStruct{
			tParsingStruct{0, "%x", &t.Sk},
			tParsingStruct{1, "%d", &t.RefCnt},
			tParsingStruct{2, "%d", &t.Type},
			tParsingStruct{3, "%x", &t.Proto},
			tParsingStruct{4, "%d", &t.Iface},
			tParsingStruct{5, "%d", &t.R},
			tParsingStruct{6, "%d", &t.Rmem},
			tParsingStruct{7, "%d", &t.User},
			tParsingStruct{8, "%d", &t.Inode},
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
