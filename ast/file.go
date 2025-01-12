package ast

type File struct {
	Filename string
	Src      []byte
}

func NewFile(filename string, src []byte) *File {
	return &File{Filename: filename, Src: src}
}
