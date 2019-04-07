package fd

type FDPipe struct {
	FD
	Address uint
}

func NewFDPipe(i uintptr, target string) *FDPipe {
	s := FDPipe{}
	s.Number = i

	return &s
}
func (f *FDPipe) GetType() string {
	return "Pipe"
}
