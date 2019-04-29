package fd

import (
	"fmt"
)

type FDFile struct {
	FD
	Path string
}

func NewFDFile(pid uint64, i uint64, target string) *FDFile {
	f := FDFile{}
	f.Number = i
	f.PID = pid
	f.Path = target

	f.SetFileInfo()

	return &f
}

func (f *FDFile) GetType() string {
	return "File"
}

func (f *FDFile) GetInfo() string {
	return fmt.Sprintf("%-8s %s", "PATH", f.Path)
}
