package fd

import (
	"strings"
)

type FDInterface interface {
	GetType() string
	GetNumber() uintptr
	GetInode() uintptr
}

type FD struct {
	Number uintptr
	Inode  uintptr
}

func NewFD(i uintptr, target string) FDInterface {
	switch {
	case strings.HasPrefix(target, "socket:"):
		return NewFDSocket(i, target)
	case strings.HasPrefix(target, "pipe:"):
		return NewFDPipe(i, target)
	case strings.HasPrefix(target, "anon_inode:"):
		return NewFDAnonInode(i, target)
	default:
		return NewFDFile(i, target)
	}
}

func (f *FD) GetNumber() uintptr {
	return f.Number
}

func (f *FD) GetInode() uintptr {
	return f.Inode
}

func (f *FD) bracketParse(s string) string {
	start := strings.LastIndex(s, "[")
	end := strings.LastIndex(s, "]")
	if start != -1 && end != -1 {
		return s[start+1 : end]
	}
	return ""
}
