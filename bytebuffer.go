package bytebufferpool

import (
	"io"
	"unicode/utf8"

	"github.com/go-asphyxia/conversion"
)

type (
	ByteBuffer struct {
		bytes  []byte
		offset int
	}
)

func NewByteBuffer(l, c int) (b *ByteBuffer) {
	b = &ByteBuffer{
		bytes: make([]byte, l, c),
	}

	return
}

func (b *ByteBuffer) Copy() (target *ByteBuffer) {
	target = &ByteBuffer{
		bytes: make([]byte, len(b.bytes), cap(b.bytes)),
	}

	copy(target.bytes, b.bytes)
	return
}

func (b *ByteBuffer) Len() (l int) {
	l = len(b.bytes)
	return
}

func (b *ByteBuffer) Cap() (c int) {
	c = cap(b.bytes)
	return
}

func (b *ByteBuffer) Grow(n int) {
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

func (b *ByteBuffer) Reset() {
	b.bytes = b.bytes[:0]
	b.offset = 0
}

func (b *ByteBuffer) Close() (err error) {
	b.bytes = nil
	return
}

func (b *ByteBuffer) Bytes() (target []byte) {
	target = b.bytes
	return
}

func (b *ByteBuffer) String() (target string) {
	target = conversion.BytesToStringNoCopy(b.bytes)
	return
}

func (b *ByteBuffer) CopyBytes() (target []byte) {
	target = make([]byte, len(b.bytes), cap(b.bytes))
	copy(target, b.bytes)
	return
}

func (b *ByteBuffer) CopyString() (target string) {
	target = string(b.bytes)
	return
}

func (b *ByteBuffer) Set(source []byte) {
	b.bytes = append(b.bytes[:0], source...)
}

func (b *ByteBuffer) SetString(source string) {
	b.bytes = append(b.bytes[:0], source...)
}

func (b *ByteBuffer) Write(source []byte) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *ByteBuffer) WriteByte(source byte) (err error) {
	b.bytes = append(b.bytes, source)
	return
}

func (b *ByteBuffer) WriteRune(source rune) (n int, err error) {
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

func (b *ByteBuffer) WriteString(source string) (n int, err error) {
	b.bytes = append(b.bytes, source...)
	n = len(source)
	return
}

func (b *ByteBuffer) ReadFrom(source io.Reader) (n int64, err error) {
	i := len(b.bytes)
	c := cap(b.bytes)

	r := 0

	for {
		if i == c {
			c = (c + 16) * 2

			temp := make([]byte, c)
			copy(temp, b.bytes)

			b.bytes = temp
		}

		r, err = source.Read(b.bytes[i:c])

		n += int64(r)
		i += r

		b.bytes = b.bytes[:i]

		if err != nil || i < c {
			if err == io.EOF {
				err = nil
			}

			return
		}
	}
}

func (b *ByteBuffer) Read(target []byte) (n int, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	n = copy(target, b.bytes[b.offset:])
	b.offset += n
	return
}

func (b *ByteBuffer) ReadByte() (target byte, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	target = b.bytes[b.offset]
	b.offset++
	return
}

func (b *ByteBuffer) ReadRune() (target rune, n int, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	target, n = utf8.DecodeRune(b.bytes[b.offset:])
	b.offset += n
	return
}

func (b *ByteBuffer) ReadString(target string) (n int, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	n = copy(conversion.StringToBytesNoCopy(target), b.bytes[b.offset:])
	b.offset += n
	return
}

func (b *ByteBuffer) WriteTo(target io.Writer) (n int64, err error) {
	if len(b.bytes) <= b.offset {
		err = io.EOF
		return
	}

	wrote, err := target.Write(b.bytes[b.offset:])
	b.offset += wrote
	n = int64(wrote)
	return
}
