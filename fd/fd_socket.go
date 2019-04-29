package fd

import (
	"strconv"
)

type FDSocket struct {
	FD
	Address uint
}

func NewFDSocket(pid uint64, i uint64, target string) *FDSocket {
	s := FDSocket{}
	s.Number = i
	s.PID = pid

	inode_str := s.bracketParse(target)
	if inode_str != "" {
		inode, err := strconv.ParseInt(inode_str, 10, 64)
		if err == nil {
			s.Inode = uint64(inode)
		}
	}
	s.SetFileInfo()
	return &s
}

func (f *FDSocket) GetType() string {
	return "Socket"
}

func (f *FDSocket) GetInfo() string {
	return ""
}
