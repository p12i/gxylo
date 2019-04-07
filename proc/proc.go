package proc

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
)

type Proc struct {
	PID int
	FDS []FD
}

func NewProc(pid int) (Proc, error) {
	procPath := path.Join("/proc", strconv.Itoa(pid))
	procInfo, err := os.Stat(procPath)
	if err != nil {
		return Proc{}, err
	}
	if !procInfo.IsDir() {
		return Proc{}, errors.New("Given path is not directory")
	}
	fmt.Println(procInfo)
	p := Proc{}
	p.PID = pid
	return p, err
}

func (p *Proc) fileDescriptorsLinkNames(path string) ([]string, error) {
	proc_dir, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	fds, err := proc_dir.Readdirnames(-1)
	if err != nil {
		return []string{}, err
	}
	proc_dir.Close()
	return fds, nil
}

func (p *Proc) ParseFDS() error {
	connections := ConnectionList{}

	if err := connections.ParseConnections(); err != nil {
		return err
	}

	var f FDInterface
	root_path := path.Join("/proc", strconv.Itoa(p.PID), "fd")
	fds, err := p.fileDescriptorsLinkNames(root_path)
	if err != nil {
		return err
	}
	p.FDS = make([]FD, len(fds))

	for _, elem := range fds {
		fd_path := path.Join(root_path, fmt.Sprint(elem))
		fd_i64, err := strconv.ParseInt(elem, 10, 32)
		fd_int := uintptr(fd_i64)
		if err != nil {
			return err
		}
		target, err := os.Readlink(fd_path)
		if err != nil {
			return err
		}

		f = NewFD(fd_int, target)
		fmt.Printf("%s -> %s\n", elem, target)
		fmt.Printf("%d -> %s [%d]\n", f.GetNumber(), f.GetType(), f.GetInode())
		if c := connections.GetConnection(f.GetInode()); c != nil {
			fmt.Println(c.String())
		}
	}
	fmt.Println(connections.String())

	return nil
}
