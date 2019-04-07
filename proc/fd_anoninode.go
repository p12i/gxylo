package proc

import (
	"fmt"
	"strings"
)

type FDAnonInode struct {
	FD
	SubType string
}

func NewFDAnonInode(i uintptr, target string) *FDAnonInode {
	s := FDAnonInode{}
	s.Number = i
	start := strings.LastIndex(target, "[")
	end := strings.LastIndex(target, "]")
	if start != -1 && end != -1 {
		s.SubType = target[start:end]
	}

	return &s
}
func (f *FDAnonInode) GetType() string {
	return fmt.Sprintf("Anon Inode[%s]", f.SubType)
}
