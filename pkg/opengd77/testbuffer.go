package opengd77

import "bytes"

type (
	TestBuffer struct {
		ReadBuffer, WriteBuffer *bytes.Buffer
	}
)

func (buf *TestBuffer) Read(dest []byte) (int, error) {
	return buf.ReadBuffer.Read(dest)
}

func (buf *TestBuffer) Write(src []byte) (int, error) {
	return buf.WriteBuffer.Write(src)
}

func (buf *TestBuffer) Close() error {
	return nil
}

func NewTestBuffer(readContent []byte) *TestBuffer {
	buf := TestBuffer{}
	buf.ReadBuffer = bytes.NewBuffer(readContent)
	buf.WriteBuffer = new(bytes.Buffer)

	return &buf
}
