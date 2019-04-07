package proc

import (
	"strconv"
)

type FDSocket struct {
	FD
	Address uint
}

func NewFDSocket(i uintptr, target string) *FDSocket {
	s := FDSocket{}
	s.Number = i

	inode_str := s.bracketParse(target)
	if inode_str != "" {
		inode, err := strconv.ParseInt(inode_str, 10, 64)
		if err == nil {
			s.Inode = uintptr(inode)
		}
	}
	return &s
}

func (f *FDSocket) GetType() string {
	return "Socket"
}
