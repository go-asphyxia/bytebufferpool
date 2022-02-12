package bytebufferpool

import (
	"io"
	"unicode/utf8"

	"github.com/go-asphyxia/conversion"
)

type (
	B struct {
		bytes  []byte
		offset int
	}
)

func Buffer(l, c int) (b *B) {
	b = &B{
		bytes: make([]byte, l, c),
	}

	return
}

func (b *B) Copy() (target *B) {
	target = &B{
		bytes: make([]byte, len(b.bytes), cap(b.bytes)),
	}

	copy(target.bytes, b.bytes)
	return
}

func (b *B) Len() (l int) {
	l = len(b.bytes)
	return
}

func (b *B) Cap() (c int) {
	c = cap(b.bytes)
	return
}

func (b *B) Grow(n int) {
	s := len(b.bytes) + n

	if s <= cap(b.bytes) {
		b.bytes = b.bytes[:s]
		return
	}

	temp := make([]byte, s, s)
	copy(temp, b.bytes)

	b.bytes = temp
	return
}

func (b *B) Reset() {
	b.bytes = b.bytes[:0]
	b.offset = 0
}

func (b *B) Close() (err error) {
	b.bytes = nil
	return
}

func (b *B) Bytes() (target []byte) {
	target = b.bytes
	return
}

func (b *B) String() (target string) {
	target = conversion.BytesToStringNoCopy(b.bytes)
	return
}

func (b *B) CopyBytes() (target []byte) {
	target = make([]byte, len(b.bytes), cap(b.bytes))
	copy(target, b.bytes)
	return
}

func (b *B) CopyString() (target string) {
	target = string(b.bytes)
	return
}

func (b *B) Set(source []byte) {
	b.bytes = append(b.bytes[:0], source...)
}

func (b *B) SetString(source string) {
	b.bytes = append(b.bytes[:0], source...)
}

func (b *B) Write(source []byte) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *B) WriteByte(source byte) (err error) {
	b.bytes = append(b.bytes, source)
	return
}

func (b *B) WriteRune(source rune) (n int, err error) {
	l := len(b.bytes)

	s := l + utf8.UTFMax

	if s <= cap(b.bytes) {
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

func (b *B) WriteString(source string) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *B) ReadFrom(source io.Reader) (n int64, err error) {
	var (
		l int
		c int
		r int
	)

f:
	for {
		l = len(b.bytes)
		c = cap(b.bytes)

		if l < c {
			r, err = source.Read(b.bytes[l:])
			n += int64(r)
			if err != nil {
				return
			}

			continue f
		}

		size := c * 2

		temp := make([]byte, size, size)
		copy(temp, b.bytes)

		r, err = source.Read(temp[l:])
		n += int64(r)
		if err != nil {
			return
		}

		b.bytes = temp
	}
}

func (b *B) Read(target []byte) (n int, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	n = copy(target, b.bytes[b.offset:])
	b.offset += n
	return
}

func (b *B) ReadByte() (target byte, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	target = b.bytes[b.offset]
	b.offset++
	return
}

func (b *B) ReadRune() (target rune, n int, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	target, n = utf8.DecodeRune(b.bytes[b.offset:])
	b.offset += n
	return
}

func (b *B) ReadString(target string) (n int, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	n = copy(conversion.StringToBytesNoCopy(target), b.bytes[b.offset:])
	b.offset += n
	return
}

func (b *B) WriteTo(target io.Writer) (n int64, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	wrote, err := target.Write(b.bytes[b.offset:])
	b.offset += wrote
	n = int64(wrote)
	return
}
