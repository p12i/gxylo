package fd

import (
	"fmt"
	"strings"
)

type FDAnonInode struct {
	FD
	SubType string
}

func NewFDAnonInode(pid uint64, i uint64, target string) *FDAnonInode {
	s := FDAnonInode{}
	s.Number = i
	s.PID = pid
	s.SubType = s.bracketParse(target)
	if len(s.SubType) == 0 && strings.Contains(target, ":") {
		start := strings.LastIndex(target, ":")
		s.SubType = target[start+1:]
	}

	s.SetFileInfo()
	return &s
}
func (f *FDAnonInode) GetType() string {
	return fmt.Sprintf("Anon Inode")
}

func (f *FDAnonInode) GetInfo() string {
	return fmt.Sprintf("%-8s %s", "SUBTYPE", f.SubType)
}
