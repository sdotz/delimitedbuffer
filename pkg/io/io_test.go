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
