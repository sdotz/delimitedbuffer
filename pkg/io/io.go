package io

import (
	"encoding/binary"
	"math"
	"github.com/pkg/errors"
	"bytes"
)

//DelimitedBuffer embeds bytes.Buffer and stores arbitrary binary data prefixed by its 4-byte size
type DelimitedBuffer struct {
	bytes.Buffer
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
	b := make([]byte, 4)
	if _, err := f.Buffer.Read(b); err != nil {
		return []byte{}, err
	}

	data := make([]byte, binary.LittleEndian.Uint32(b))
	if _, err := f.Buffer.Read(data); err != nil {
		return []byte{}, err
	}

	return data, nil
}
