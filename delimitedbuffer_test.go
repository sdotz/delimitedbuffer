package delimitedbuffer

import (
	"testing"
	"bytes"
	"compress/gzip"
	"io"
)

func TestShouldWriteMultipleDelimitedBinaries(t *testing.T) {
	b1 := []byte("foo")
	b2 := []byte("2😗↗️₯℃×≸∛23rahhh soooooo😍")

	dbuf := DelimitedBuffer{}

	dbuf.Write(b1)
	dbuf.Write(b2)

	bb1, err := dbuf.ReadNext()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b1, bb1) {
		t.Error("did not match b1")
	}

	bb2, err := dbuf.ReadNext()
	if !bytes.Equal(b2, bb2) {
		t.Error("did not match b2")
	}
}

func TestShouldReturnNoBytesAtEndOfChunk(t *testing.T) {
	b1 := []byte("foo")
	b2 := []byte("2😗↗️₯℃×≸∛23rahhh soooooo😍")
	dbuf := DelimitedBuffer{}
	dbuf.Write(b1)
	dbuf.Write(b2)
	p := make([]byte, 20)

	n, _ := dbuf.Read(p)
	if string(p[:n]) != "foo" {
		t.Error("read data did not match expected")
	}

	n, _ = dbuf.Read(p)
	if n != 0 {
		t.Error("read data should have returned 0 bytes at end of chunk")
	}
}

func TestReadWritePassthrough(t *testing.T) {
	b1 := []byte("123")
	b2 := []byte("2😗↗️₯℃×≸∛23raddddhhh soooooo😍")

	dbuf := DelimitedBuffer{}

	zw := gzip.NewWriter(&dbuf)

	zw.Write(b1)
	zw.Flush()
	zw.Write(b2)
	zw.Flush()

	if err := zw.Close(); err != nil {
		panic(err)
	}

	gzr, _ := gzip.NewReader(&dbuf)

	var output [][]byte
	var message []byte
	for {
		readBytes := make([]byte, 4)
		n, err := gzr.Read(readBytes)
		message = append(message, readBytes[:n]...)
		if n < len(readBytes) {
			output = append(output, message)
			message = []byte{}
		}

		if err == io.EOF || n == 0 {
			break
		}
	}
	if !bytes.Equal(b1, output[0]) {
		t.Error("byte slices did not match for b1")
	}

	if !bytes.Equal(b2, output[1]) {
		t.Error("byte slices did not match for b2")
	}
}
