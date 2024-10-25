package utils

import (
	"bufio"
	"errors"
	"io"
)

// DEFAULT_BUFFER_SIZE Default buffer size to use. Let's keep it not too small.
const DEFAULT_BUFFER_SIZE = 256 * 1024 // 256 kB buffer

type bufReadCloser struct {
	closed bool
	read   *bufio.Reader
	closer io.Closer
}

// NewBufReadCloser creates a new buffered reader that implements io.ReadCloser.
// Closing this buffered reader, will also close the underlying reader after flushing
// remaining data from the buffer.
func NewBufReadCloser(reader io.ReadCloser) io.ReadCloser {
	return &bufReadCloser{
		read:   bufio.NewReaderSize(reader, DEFAULT_BUFFER_SIZE),
		closer: reader,
	}
}

func (buf *bufReadCloser) Read(b []byte) (int, error) {
	if buf.closed {
		return 0, io.ErrClosedPipe
	}
	return buf.read.Read(b)
}

func (buf *bufReadCloser) Close() error {
	buf.closed = true
	return buf.closer.Close()
}

type bufWriteCloser struct {
	closed bool
	write  *bufio.Writer
	closer io.Closer
}

// NewBufWriteCloser creates a new buffered reader that implements io.WriteCloser.
// Closing this buffered writer, will flush the buffer and close the underlying writer.
func NewBufWriteCloser(writer io.WriteCloser) io.WriteCloser {
	return &bufWriteCloser{
		write:  bufio.NewWriterSize(writer, DEFAULT_BUFFER_SIZE),
		closer: writer,
	}
}

func (buf *bufWriteCloser) Write(b []byte) (int, error) {
	if buf.closed {
		return 0, io.ErrClosedPipe
	}
	return buf.write.Write(b)
}

func (buf *bufWriteCloser) Close() error {
	buf.closed = true
	return errors.Join(
		buf.write.Flush(),
		buf.closer.Close(),
	)
}

// ReadAllMax reads from r until an error or EOF or max bytes is reached and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported. If the data are larger to max bytes, and if max is > 0
// an error will be returned.
func ReadAllMax(r io.Reader, maxBytes int) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
		if maxBytes > 0 && len(b) > maxBytes {
			return b, errors.New("Input data is too large")
		}

		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
	}
}
