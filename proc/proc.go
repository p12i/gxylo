package proc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/p12i/gxylo/connections"
	"github.com/p12i/gxylo/fd"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

type Proc struct {
	PID      uint64
	FDS      []fd.FDInterface
	Exe      string
	Cwd      string
	Cmdline  string
	Environ  map[string]string
	procPath string
}

func prettyHeader(name string) string {
	return fmt.Sprintf("\n--====::[ %s ]::====--\n", strings.ToUpper(name))
}

func NewProc(pid uint64) (Proc, error) {
	procPath := path.Join("/proc", strconv.Itoa(int(pid)))
	procInfo, err := os.Stat(procPath)
	if err != nil {
		return Proc{}, err
	}
	if !procInfo.IsDir() {
		return Proc{}, errors.New("Given path is not directory")
	}

	p := Proc{}
	p.PID = pid
	p.procPath = procPath
	p.Exe, _ = p.procFsReadlink("exe")
	p.Cwd, _ = p.procFsReadlink("cwd")
	p.parseCmdline()
	p.setupEnviron()
	return p, nil
}
func (p *Proc) parseCmdline() {
	var value []string
	cmd_fp, _ := p.procFsFileOpen("cmdline")
	scanner := bufio.NewScanner(cmd_fp)
	scanner.Split(scanBytes)
	for scanner.Scan() {
		value = append(value, scanner.Text())

	}
	cmd_fp.Close()
	p.Cmdline = strings.Join(value, " ")
}

func scanBytes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, 0x0); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func (p *Proc) setupEnviron() {
	fp, _ := p.procFsFileOpen("environ")
	p.Environ = make(map[string]string)
	scanner := bufio.NewScanner(fp)
	scanner.Split(scanBytes)
	for scanner.Scan() {
		data := scanner.Text()
		if idx := strings.Index(data, "="); idx >= 0 {
			p.Environ[data[0:idx]] = data[idx+1:]
		}
	}
	fp.Close()
}

func (p *Proc) getProcFsPath() string {
	return p.procPath
}

func (p *Proc) procFsReadlink(file ...string) (string, error) {
	return os.Readlink(p.getProcFsFilePath(file...))
}

func (p *Proc) procFsFileOpen(file ...string) (*os.File, error) {
	return os.Open(p.getProcFsFilePath(file...))
}

func (p *Proc) getProcFsFilePath(file ...string) string {
	args := []string{p.procPath}
	args = append(args, file...)
	return path.Join(args...)
}

func (p *Proc) procFsReadContent(filename string, w io.Writer) {
	fp, _ := p.procFsFileOpen(filename)
	defer fp.Close()
	reader := bufio.NewReader(fp)
	reader.WriteTo(w)
}

func (p *Proc) Info(connections_list *connections.ConnectionList) string {
	var value strings.Builder

	value.WriteString(prettyHeader("command"))
	value.WriteString(fmt.Sprintf("Exe:     %s\nCwd:     %s\nCmdline: %s\n", p.Exe, p.Cwd, p.Cmdline))

	value.WriteString(prettyHeader("environ"))
	for k, v := range p.Environ {
		value.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}

	value.WriteString(prettyHeader("status"))
	p.procFsReadContent("status", &value)

	value.WriteString(prettyHeader("io"))
	p.procFsReadContent("io", &value)

	value.WriteString(prettyHeader("file descriptors"))
	for _, elem := range p.FDS {
		value.WriteString(fd.FDToString(elem))
		inode := elem.GetInode()
		if c, ok := connections_list.GetConnection(inode); ok && inode != 0 {
			value.WriteString(c.String())
		} else if elem.GetType() == "Socket" {
			// FIXME STAT
			value.WriteString(fmt.Sprintf("%s\n", elem.GetFileStat_t()))
		} else {
			value.WriteString("\n")
		}
	}

	return value.String()
}

func (p *Proc) fileDescriptorsLinkNames() ([]string, error) {
	if proc_dir, err := p.procFsFileOpen("fd"); err != nil {
		return []string{}, err
	} else {
		defer proc_dir.Close()
		if fds, err := proc_dir.Readdirnames(-1); err != nil {
			return []string{}, err
		} else {
			return fds, nil
		}
	}
}

func (p *Proc) ParseFDS() error {
	var fd_i64 int64
	fds, err := p.fileDescriptorsLinkNames()
	if err != nil {
		return err
	}
	p.FDS = make([]fd.FDInterface, len(fds))

	for i, elem := range fds {
		if fd_i64, err = strconv.ParseInt(elem, 10, 64); err != nil {
			return err
		}
		target, err := p.procFsReadlink("fd", fmt.Sprint(elem))
		if err != nil {
			return err
		}

		p.FDS[i] = fd.NewFD(p.PID, uint64(fd_i64), target)

	}

	return nil
}
