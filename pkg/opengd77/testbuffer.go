package opengd77

import "bytes"

type (
	// A TestBuffer is using for testing code that expects to
	// read/write to/from a device or file. TestBuffer implements
	// the io.ReadWriteCloser interface.
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
