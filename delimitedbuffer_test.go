package delimitedbuffer

import (
	"bytes"
	"testing"
)

func TestProtoBufWriterAndReader(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := NewProtoBufWriter(buf)

	// Write some data
	data := []byte("Hello, World!")
	_, err := writer.Write(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read it back
	reader := NewProtoBufReader(buf)
	readData, err := reader.Read()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(data, readData) {
		t.Fatalf("expected %v, got %v", data, readData)
	}
}

func TestProtoBufReaderWithNoData(t *testing.T) {
	buf := new(bytes.Buffer)
	reader := NewProtoBufReader(buf)

	// Try to read from an empty buffer
	_, err := reader.Read()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestProtoBufWriterWithNilData(t *testing.T) {
	buf := new(bytes.Buffer)
	writer := NewProtoBufWriter(buf)

	// Try to write nil data
	_, err := writer.Write(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
