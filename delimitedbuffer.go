package delimitedbuffer

import (
	"encoding/binary"
	"io"
)

type ProtoBufWriter struct {
	w io.Writer
}

func NewProtoBufWriter(w io.Writer) *ProtoBufWriter {
	return &ProtoBufWriter{w: w}
}

func (p *ProtoBufWriter) Write(data []byte) (int, error) {
	length := uint32(len(data))
	err := binary.Write(p.w, binary.LittleEndian, length)
	if err != nil {
		return 0, err
	}
	return p.w.Write(data)
}

type ProtoBufReader struct {
	r io.Reader
}

func NewProtoBufReader(r io.Reader) *ProtoBufReader {
	return &ProtoBufReader{r: r}
}

func (p *ProtoBufReader) Read() ([]byte, error) {
	var length uint32
	err := binary.Read(p.r, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	data := make([]byte, length)
	_, err = io.ReadFull(p.r, data)
	return data, err
}
