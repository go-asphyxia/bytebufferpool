package bytebufferpool

import (
	"io"
	"unicode/utf8"

	"github.com/go-asphyxia/conversion"
)

type (
	Buffer struct {
		bytes []byte
	}
)

func (b *Buffer) Len() (length int) {
	length = len(b.bytes)
	return
}

func (b *Buffer) Cap() (capacity int) {
	capacity = cap(b.bytes)
	return
}

func (b *Buffer) Grow(n, preallocate int) {
	l := len(b.bytes)
	c := cap(b.bytes)

	s := l + n

	if s <= c && preallocate <= 0 {
		b.bytes = b.bytes[:s]
		return
	}

	temp := make([]byte, s, (s + preallocate))
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
	l := len(b.bytes)
	c := cap(b.bytes)

	s := l + utf8.UTFMax

	if s <= c {
		b.bytes = b.bytes[:s]

		n = utf8.EncodeRune(b.bytes[l:], source)
		return
	}

	temp := make([]byte, s, s)
	copy(temp, b.bytes)

	n = utf8.EncodeRune(temp[l:], source)

	b.bytes = temp
	return
}

func (b *Buffer) WriteString(source string) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *Buffer) ReadFrom(source io.Reader) (n int64, err error) {
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
