package delimitedbuffer

import (
	"encoding/binary"
	"math"
	"github.com/pkg/errors"
	"bytes"
	"io"
)

//DelimitedBuffer embeds bytes.Buffer and stores arbitrary binary data prefixed by its 4-byte size
type DelimitedBuffer struct {
	bytes.Buffer
	remainingChunkBytes uint32
	hasBeenRead         bool
}

var ErrEndOfDatum = errors.New("Reached the end of the datum, but there are more")

func NewDelimitedBuffer(buf []byte) *DelimitedBuffer {
	return &DelimitedBuffer{Buffer: *bytes.NewBuffer(buf)}
}

// Write takes a byte slice and writes it to the DelimitedBuffer
func (f *DelimitedBuffer) Write(data []byte) (int, error) {
	bLen := len(data)
	if bLen > math.MaxUint32 {
		return 0, errors.Errorf("data size: %d, exceeded max that can be expressed in 4 bytes: %d", bLen, math.MaxUint32)
	}
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(bLen))
	wrote1, err := f.Buffer.Write(b)
	if err != nil {
		return 0, err
	}
	wrote2, err := f.Buffer.Write(data)
	if err != nil {
		return wrote1, err
	}
	return wrote1 + wrote2, nil
}

//ReadNext reads the next byte slice from the buffer and returns it
func (f *DelimitedBuffer) ReadNext() ([]byte, error) {
	bSize, err := f.getNextChunkSize()
	if err != nil {
		return []byte{}, err
	}
	data := make([]byte, bSize)
	if _, err := f.Buffer.Read(data); err != nil {
		return []byte{}, err
	}

	return data, nil
}

func (f *DelimitedBuffer) getNextChunkSize() (uint32, error) {
	if f.remainingChunkBytes > 0 {
		return 0, errors.New("cannot get next chunk size since previous chunk has not been completely read")
	}
	b := make([]byte, 4)
	n, err := f.Buffer.Read(b)
	if err != nil {
		return 0, err
	}
	if n < len(b) {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.LittleEndian.Uint32(b), nil
}

func (f *DelimitedBuffer) ReadByte() (byte, error) {
	if f.remainingChunkBytes == 0 {
		nextChunkSize, err := f.getNextChunkSize()
		if err != nil {
			return 0, err
		}
		f.remainingChunkBytes = nextChunkSize
	}
	b, err := f.Buffer.ReadByte()
	if err != nil {
		return 0, err
	}
	f.remainingChunkBytes = f.remainingChunkBytes - 1
	return b, nil
}

// Read reads up to the end of the next chunk
func (f *DelimitedBuffer) Read(b []byte) (int, error) {
	if f.remainingChunkBytes == 0 {
		byteSize, err := f.getNextChunkSize()
		if err != nil {
			// we have read that last datum, this is EOF
			return 0, err
		}
		f.remainingChunkBytes = byteSize
		if f.hasBeenRead {
			// we have read to the end of the datum but there are more
			return 0, nil
		}
	}

	// read at most f.remainingChunkBytes
	var buf []byte
	if uint32(len(b)) > f.remainingChunkBytes {
		buf = make([]byte, f.remainingChunkBytes)
	} else {
		buf = make([]byte, len(b))
	}

	n, err := f.Buffer.Read(buf)
	if err != nil {
		return n, err
	}
	f.hasBeenRead = true
	f.remainingChunkBytes -= uint32(n)

	copy(b, buf)

	return n, nil
}
