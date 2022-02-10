package bytebufferpool

import (
	"errors"
	"io"
	"unicode/utf8"

	"github.com/go-asphyxia/conversion"
)

type (
	Buffer struct {
		bytes []byte
	}
)

const (
	preallocate int = 16
)

func (b *Buffer) Length() (length int) {
	length = len(b.bytes)
	return
}

func (b *Buffer) Capacity() (capacity int) {
	capacity = cap(b.bytes)
	return
}

func (b *Buffer) Preallocate(n int) {
	l := len(b.bytes)
	c := cap(b.bytes)

	temp := make([]byte, l, (c + n))
	copy(temp, b.bytes)

	b.bytes = temp
}

func (b *Buffer) Grow(n int) {
	l := len(b.bytes)
	c := cap(b.bytes)

	s := l + n

	if s <= c {
		b.bytes = b.bytes[:s]
		return
	}

	temp := make([]byte, s, (s + preallocate))
	copy(temp, b.bytes)

	b.bytes = temp
}

func (b *Buffer) PreallocateAndGrow(n, grow int) {
	l := len(b.bytes)
	c := cap(b.bytes)

	temp := make([]byte, (l + grow), (c + n + grow))
	copy(temp, b.bytes)

	b.bytes = temp
}

func (b *Buffer) Reset() {
	b.bytes = b.bytes[:0]
}

func (b *Buffer) Close() (err error) {
	b.bytes = nil
	return
}

func (b *Buffer) Bytes() (target []byte) {
	target = b.bytes
	return
}

func (b *Buffer) String() (target string) {
	target = conversion.BytesToStringNoCopy(b.bytes)
	return
}

func (b *Buffer) CopyBytes() (target []byte) {
	l := len(b.bytes)
	c := cap(b.bytes)

	target = make([]byte, l, c)
	copy(target, b.bytes)
	return
}

func (b *Buffer) CopyString() (target string) {
	target = string(b.bytes)
	return
}

func (b *Buffer) Set(source []byte) {
	b.bytes = append(b.bytes[:0], source...)
}

func (b *Buffer) SetString(source string) {
	b.bytes = append(b.bytes[:0], source...)
}

func (b *Buffer) Write(source []byte) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *Buffer) WriteByte(source byte) (err error) {
	b.bytes = append(b.bytes, source)
	return
}

func (b *Buffer) WriteRune(source rune) (n int, err error) {
	if uint32(source) < utf8.RuneSelf {
		b.WriteByte(byte(source))
		n = 1
		return
	}

	l := len(b.bytes)

	b.Grow(utf8.UTFMax)

	n = utf8.EncodeRune(b.bytes[l:], source)

	b.bytes = b.bytes[:l+n]
	return
}

func (b *Buffer) WriteString(source string) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *Buffer) ReadFrom(source io.Reader) (n int64, err error) {
	eof := false
	read := 0

	for !eof {
		read, err = source.Read(b.bytes[:read])
		eof = errors.Is(err, io.EOF)
		if err != nil && !eof {
			return
		}

		n += int64(read)
		b.Grow(read * 2)
	}

	return
}

func (b *Buffer) Read(target []byte) (n int, err error) {
	n = copy(target, b.bytes)
	return
}

func (b *Buffer) ReadString(target string) (n int, err error) {
	n = copy(conversion.StringToBytesNoCopy(target), b.bytes)
	return
}

func (b *Buffer) WriteTo(target io.Writer) (n int64, err error) {
	wrote, err := target.Write(b.bytes)
	n = int64(wrote)
	return
}
