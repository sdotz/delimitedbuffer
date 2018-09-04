package manio

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

func TestReadWritePassthrough(t *testing.T){
	b2 := []byte("2😗↗️₯℃×≸∛23rahhh soooooo😍")

	dbuf := DelimitedBuffer{}

	zw := gzip.NewWriter(&dbuf)

	zw.Write(b2)

	if err := zw.Close(); err != nil {
		panic(err)
	}

	zr, _:= gzip.NewReader(&dbuf)

	outbuf := bytes.NewBuffer([]byte{})

	if _, err := io.Copy(outbuf, zr); err != nil {
		panic(err)
	}

	if err := zr.Close(); err != nil {
		panic(err)
	}

	if  !bytes.Equal(outbuf.Bytes(), b2) {
		t.Error("output bytes did not match expected output")
	}
}
