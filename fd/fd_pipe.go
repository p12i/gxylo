package fd

import (
	"fmt"
)

type FDPipe struct {
	FD
}

func NewFDPipe(pid uint64, i uint64, target string) *FDPipe {
	s := FDPipe{}
	s.Number = i
	s.PID = pid
	fmt.Sscanf(s.bracketParse(target), "%d", &s.Inode)

	s.SetFileInfo()
	return &s
}
func (f *FDPipe) GetType() string {
	return "Pipe"
}
func (f *FDPipe) GetInfo() string {
	return ""
}
