package bytebufferpool

import (
	"unicode/utf8"

	"github.com/go-asphyxia/conversion"
)

type (
	Buffer struct {
		bytes []byte
	}
)

func (b *Buffer) Length() (length int) {
	length = len(b.bytes)
	return
}

func (b *Buffer) Capacity() (capacity int) {
	capacity = cap(b.bytes)
	return
}

func (b *Buffer) Grow(n, preallocate int) {
	l := len(b.bytes)
	c := cap(b.bytes)

	s := l + n

	if (c-l) > n && preallocate < 1 {
		b.bytes = b.bytes[:s]
		return
	}

	new := make([]byte, s, (s + preallocate))

	for i := 0; i < l; i++ {
		new[i] = b.bytes[i]
	}

	b.bytes = new
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
	target = conversion.BytesToString(b.bytes)
	return
}

func (b *Buffer) Copy() (target []byte) {
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

	b.Grow(utf8.UTFMax, 0)

	n = utf8.EncodeRune(b.bytes[l:], source)
	return
}

func (b *Buffer) WriteString(source string) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *Buffer) Read() {

}
