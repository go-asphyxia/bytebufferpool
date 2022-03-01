package bytebufferpool

import (
	"io"
	"unicode/utf8"

	"github.com/go-asphyxia/conversion"
)

type (
	ByteBuffer struct {
		Bytes  []byte
		Offset int
	}
)

func NewByteBuffer() (b *ByteBuffer) {
	b = &ByteBuffer{
		Bytes: make([]byte, 0),
	}

	return
}

func (b *ByteBuffer) Copy() (target *ByteBuffer) {
	target = &ByteBuffer{
		Bytes: make([]byte, len(b.Bytes), cap(b.Bytes)),
	}

	copy(target.Bytes, b.Bytes)
	return
}

func (b *ByteBuffer) Len() (l int) {
	l = len(b.Bytes)
	return
}

func (b *ByteBuffer) Cap() (c int) {
	c = cap(b.Bytes)
	return
}

func (b *ByteBuffer) Grow(n int) {
	s := len(b.Bytes) + n

	if s <= cap(b.Bytes) {
		b.Bytes = b.Bytes[:s]
		return
	}

	temp := make([]byte, s, s)
	copy(temp, b.Bytes)

	b.Bytes = temp
	return
}

func (b *ByteBuffer) Reset() {
	b.Bytes = b.Bytes[:0]
	b.Offset = 0
}

func (b *ByteBuffer) Close() (err error) {
	b.Bytes = nil
	return
}

func (b *ByteBuffer) String() (target string) {
	target = conversion.BytesToStringNoCopy(b.Bytes)
	return
}

func (b *ByteBuffer) CopyBytes() (target []byte) {
	target = make([]byte, len(b.Bytes), cap(b.Bytes))
	copy(target, b.Bytes)
	return
}

func (b *ByteBuffer) CopyString() (target string) {
	target = string(b.Bytes)
	return
}

func (b *ByteBuffer) Set(source []byte) {
	b.Bytes = append(b.Bytes[:0], source...)
}

func (b *ByteBuffer) SetString(source string) {
	b.Bytes = append(b.Bytes[:0], source...)
}

func (b *ByteBuffer) Write(source []byte) (n int, err error) {
	b.Bytes = append(b.Bytes, source...)
	n = len(source)
	return
}

func (b *ByteBuffer) WriteByte(source byte) (err error) {
	b.Bytes = append(b.Bytes, source)
	return
}

func (b *ByteBuffer) WriteRune(source rune) (n int, err error) {
	l := len(b.Bytes)

	s := l + utf8.UTFMax

	if s <= cap(b.Bytes) {
		b.Bytes = b.Bytes[:s]

		n = utf8.EncodeRune(b.Bytes[l:], source)
		return
	}

	temp := make([]byte, s, s)
	copy(temp, b.Bytes)

	n = utf8.EncodeRune(temp[l:], source)

	b.Bytes = temp
	return
}

func (b *ByteBuffer) WriteString(source string) (n int, err error) {
	b.Bytes = append(b.Bytes, source...)
	n = len(source)
	return
}

func (b *ByteBuffer) ReadFrom(source io.Reader) (n int64, err error) {
	i := len(b.Bytes)
	c := cap(b.Bytes)

	r := 0

	for {
		if i == c {
			c = (c + 16) * 2

			temp := make([]byte, c)
			copy(temp, b.Bytes)

			b.Bytes = temp
		}

		r, err = source.Read(b.Bytes[i:c])

		n += int64(r)
		i += r

		b.Bytes = b.Bytes[:i]

		if err != nil || i < c {
			if err == io.EOF {
				err = nil
			}

			return
		}
	}
}

func (b *ByteBuffer) Read(target []byte) (n int, err error) {
	if len(b.Bytes) <= b.Offset {
		err = io.EOF
		return
	}

	n = copy(target, b.Bytes[b.Offset:])
	b.Offset += n
	return
}

func (b *ByteBuffer) ReadByte() (target byte, err error) {
	if len(b.Bytes) <= b.Offset {
		err = io.EOF
		return
	}

	target = b.Bytes[b.Offset]
	b.Offset++
	return
}

func (b *ByteBuffer) ReadRune() (target rune, n int, err error) {
	if len(b.Bytes) <= b.Offset {
		err = io.EOF
		return
	}

	target, n = utf8.DecodeRune(b.Bytes[b.Offset:])
	b.Offset += n
	return
}

func (b *ByteBuffer) ReadString(target string) (n int, err error) {
	if len(b.Bytes) <= b.Offset {
		err = io.EOF
		return
	}

	n = copy(conversion.StringToBytesNoCopy(target), b.Bytes[b.Offset:])
	b.Offset += n
	return
}

func (b *ByteBuffer) WriteTo(target io.Writer) (n int64, err error) {
	if len(b.Bytes) <= b.Offset {
		err = io.EOF
		return
	}

	wrote, err := target.Write(b.Bytes[b.Offset:])
	b.Offset += wrote
	n = int64(wrote)
	return
}
