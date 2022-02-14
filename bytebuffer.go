package bytebufferpool

import (
	"errors"
	"io"
	"unicode/utf8"

	"github.com/go-asphyxia/conversion"
)

type (
	C struct {
		Len int
		Cap int
	}

	B struct {
		bytes  []byte
		offset int
	}
)

func Buffer(configuration *C) (b *B) {
	l := 0
	c := 0

	if configuration != nil {
		l = configuration.Len
		c = configuration.Cap
	}

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
	i := len(b.bytes)
	r := 0

	c := cap(b.bytes)

	if c == 0 {
		b.bytes = make([]byte, 64)
	}

	for {
		r, err = source.Read(b.bytes[i:c])
		n += int64(r)
		i += r

		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}

			b.bytes = b.bytes[:n]
			return
		}

		c = (c + 2) * 2

		temp := make([]byte, c)
		copy(temp, b.bytes)

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
