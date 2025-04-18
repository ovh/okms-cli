package flagsmgmt

import (
	"errors"
	"io"
	"os"
	"strings"
	"sync/atomic"

	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/ovh/okms-cli/internal/utils"
)

var (
	stdinUsed  atomic.Bool
	stdoutUsed atomic.Bool
)

// StringFromArg gets the bytes passed as a CLI argument into value.
//
// - If the value starts with '@', then the string is read from the file following @. For example @~/myfile.txt.
// - If the value is equal to '-', then the string is read from stdin.
// - Otherwise, value is converted into a string and returned as is.
func StringFromArg(value string, maxBytes int) string {
	return string(BytesFromArg(value, maxBytes))
}

// BytesFromArg gets the bytes passed as a CLI argument into value.
//
// - If the value starts with '@', then the bytes are read from the file following @. For example @~/myfile.txt.
// - If the value is equal to '-', then the bytes are read from stdin.
// - Otherwise, value is converted into byte slice and returned as is.
func BytesFromArg(value string, maxBytes int) []byte {
	reader := ReaderFromArg(value)
	defer reader.Close()
	return exit.OnErr2(utils.ReadAllMax(reader, maxBytes))
}

// ReaderFromArg gets a reader to read bytes passed as CLI argument.
//
// - If the value starts with '@', then the bytes are read from the file following @. For example @~/myfile.txt.
// - If the value is equal to '-', then the bytes are read from stdin.
// - Otherwise, value is converted into byte reader and returned as is.
//
// The function returns an [io.ReadCloser]. It's the caller's responsibility to
// call the Close() method to dispose the reader and free allocated resources.
func ReaderFromArg(value string) io.ReadCloser {
	r, _ := ReaderFromArgWithSize(value)
	return r
}

// ReaderFromArgWithSize The function returns an [io.ReadCloser] and enventually the size of data to read, or -1.
// It's the caller's responsibility to call the Close() method to dispose the reader
// and free allocated resources.
func ReaderFromArgWithSize(value string) (read io.ReadCloser, size int64) {
	read, size, err := readerFromArg(value)
	exit.OnErr(err)
	return read, size
}

func readerFromArg(value string) (io.ReadCloser, int64, error) {
	switch {
	case value[0] == '@':
		value, err := utils.ExpandTilde(value[1:])
		if err != nil {
			return nil, 0, err
		}
		f, err := os.Open(value)
		if err != nil {
			return nil, 0, err
		}
		info, err := f.Stat()
		if err != nil {
			return nil, 0, err
		}
		return utils.NewBufReadCloser(f), info.Size(), nil
	case value == "-":
		if stdinUsed.Swap(true) {
			return nil, 0, errors.New("Cannot read stdin more than once")
		}
		return utils.NewBufReadCloser(os.Stdin), -1, nil
	case strings.HasPrefix(value, "\\@"):
		value = value[1:]
		fallthrough
	default:
		return io.NopCloser(strings.NewReader(value)), int64(len(value)), nil
	}
}

func WriterFromArg(value string) io.WriteCloser {
	return exit.OnErr2(writerFromArg(value))
}

func writerFromArg(value string) (io.WriteCloser, error) {
	switch value {
	case "-":
		if stdoutUsed.Swap(true) {
			return nil, errors.New("Cannot write to stdout more than once")
		}
		return os.Stdout, nil
	default:
		value, err := utils.ExpandTilde(value)
		if err != nil {
			return nil, err
		}
		f, err := os.Create(value)
		if err != nil {
			return nil, err
		}
		return utils.NewBufWriteCloser(f), nil
	}
}

// type MaybeStdin struct {
// 	value string
// }

// func (e *MaybeStdin) Set(v string) error {
// 	e.value = v
// 	return nil
// }

// func (e *MaybeStdin) Type() string {
// 	return "value|@path|-"
// }

// func (e *MaybeStdin) Reader() io.ReadCloser {
// 	if e.value == "" {
// 		return nil
// 	}
// 	return ReaderFromArg(e.value)
// }

// func (e *MaybeStdin) Bytes() []byte {
// 	r := e.Reader()
// 	if r == nil {
// 		return []byte{}
// 	}
// 	defer r.Close()
// 	return exit.OnErr2(io.ReadAll(r))
// }

// func (e *MaybeStdin) String() string {
// 	return string(e.Bytes())
// }
