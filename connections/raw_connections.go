package connections

import ()

type RawConnection struct {
	SocketConnection
}

func (t RawConnection) String() string {
	return t.SocketConnection.String("RAW")
}

//  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode
//   0: 0103000A:0035 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 40116 1 0000000000000000 100 0 0 10 0

func (l *ConnectionList) ParseRawConnections() error {
	ch := make(chan SocketConnection)
	ctrl_ch := make(chan ChannelControl)

	go parseSocketConnetions("raw", ch, ctrl_ch)
	for {
		select {
		case i := <-ch:
			l.Connections[i.Inode] = RawConnection{i}
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
