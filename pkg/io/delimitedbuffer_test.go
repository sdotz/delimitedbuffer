package io

import (
	"testing"
	"bytes"
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
