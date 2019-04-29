package fd

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
)

type FDInterface interface {
	GetType() string
	GetNumber() uint64
	GetInode() uint64
	GetInfo() string
	GetFDPath() string
	SetFileInfo() error
	GetFileStat_t() *syscall.Stat_t
}

type FD struct {
	FileInfo os.FileInfo
	Number   uint64
	Inode    uint64
	PID      uint64
}

func NewFD(pid uint64, i uint64, target string) (f FDInterface) {
	switch {
	case strings.HasPrefix(target, "socket:"):
		f = NewFDSocket(pid, i, target)
	case strings.HasPrefix(target, "pipe:"):
		f = NewFDPipe(pid, i, target)
	case strings.HasPrefix(target, "anon_inode:"):
		f = NewFDAnonInode(pid, i, target)
	default:
		f = NewFDFile(pid, i, target)
	}
	return f
}

func (f *FD) GetFDPath() string {
	pt := path.Join("/proc", strconv.Itoa(int(f.PID)), "fd", strconv.Itoa(int(f.Number)))
	return pt
}

func (f *FD) GetFileStat_t() *syscall.Stat_t {
	return f.FileInfo.Sys().(*syscall.Stat_t)
}

func (f *FD) GetUid() uint32 {
	return f.GetFileStat_t().Uid
}

func (f *FD) GetGid() uint32 {
	return f.GetFileStat_t().Gid
}

func (f *FD) GetStatIno() uint64 {
	return f.GetFileStat_t().Ino
}

func (f *FD) SetFileInfo() error {
	var f_info os.FileInfo
	var err error
	if f_info, err = os.Lstat(f.GetFDPath()); err != nil {
		return err
	}

	f.FileInfo = f_info
	if f.Inode == 0 {
		f.Inode = f.GetStatIno()
	}
	return nil
}

func (f *FD) GetNumber() uint64 {
	return f.Number
}

func (f *FD) GetInode() uint64 {
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

func FDToString(f FDInterface) string {
	return fmt.Sprintf("[%4d] %-12s - Inode: %8d| %s", f.GetNumber(), f.GetType(), f.GetInode(), f.GetInfo())

}
