package proc

type FDFile struct {
	FD
}

func NewFDFile(i uintptr, target string) *FDFile {
	f := FDFile{}
	f.Number = i

	return &f
}

func (f *FDFile) GetType() string {
	return "File"
}
