package token

// File is a source file.
type File struct {
	Filename string
	Src      []byte
}

// NewFile creates a new file
func NewFile(filename string, src []byte) *File {
	return &File{Filename: filename, Src: src}
}

// TokenBytes returns a cheap byte slice of the token text
func (f *File) TokenBytes(t Token) []byte {
	return f.Src[t.Offset():t.EndOfToken(f.Src)]
}

// TokenString allocates a new string copy of the token text
func (f *File) TokenString(t Token) string {
	return string(f.TokenBytes(t))
}

// PositionOf returns the line and column of the given token
func (f *File) PositionOf(tok Token) (line int, col int) {
	offset := tok.Offset()
	var lineoffset int
	for i, ch := range f.Src {
		if i >= offset {
			break
		}
		if ch == '\n' {
			line++
			lineoffset = i
		}
	}
	col = offset - lineoffset
	return line + 1, col + 1
}

// String returns the filename.
func (f *File) String() string {
	return f.Filename
}
