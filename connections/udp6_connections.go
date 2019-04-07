package connections

import ()

type UDP6Connection struct {
	SocketConnection
}

func (t UDP6Connection) String() string {
	return t.SocketConnection.String("UDP6")
}

//   sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode ref pointer drops
//   76: 0103000A:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000     0        0 41260 2 0000000000000000 0
//      76: 3500007F:0035 00000000:0000 07 00000000:00000000 00:00000000 00000000   101        0 26015 2 0000000000000000 0
//
func (l *ConnectionList) ParseUDP6Connections() error {
	ch := make(chan SocketConnection)
	ctrl_ch := make(chan ChannelControl)

	go parseSocketConnetions("udp6", ch, ctrl_ch)
	for {
		select {
		case i := <-ch:
			l.Connections[i.Inode] = UDP6Connection{i}
		case i := <-ctrl_ch:
			if i.MsgType == CH_CTRL_ERR {
				return i.Error
			} else if i.MsgType == CH_CTRL_QUIT {
				return nil
			}
		}
	}

	return nil
}
