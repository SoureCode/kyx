package shell

import (
	"bytes"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

type OutputWriter struct {
	lock    sync.Mutex
	buffer  bytes.Buffer
	lineBuf bytes.Buffer
	Lines   chan string
}

func NewOutputWriter() *OutputWriter {
	ow := &OutputWriter{
		Lines: make(chan string, 100), // buffer to avoid blocking
	}

	return ow
}

func (ow *OutputWriter) Write(p []byte) (int, error) {
	ow.lock.Lock()
	defer ow.lock.Unlock()

	n, err := ow.buffer.Write(p)
	if err != nil {
		return n, errors.Wrap(err, "failed to write to buffer")
	}
	// Mirror to lineBuf for line splitting
	_, _ = ow.lineBuf.Write(p)

	// Stream out complete lines
	for {
		line, err := ow.lineBuf.ReadString('\n')
		if err == io.EOF {
			ow.lineBuf = *bytes.NewBufferString(line)
			break
		}
		// Remove trailing '\n' (and '\r' for Windows CRLF)
		line = strings.TrimRight(line, "\r\n")
		ow.Lines <- line
	}

	return n, nil
}

func (ow *OutputWriter) Buffer() []byte {
	ow.lock.Lock()
	defer ow.lock.Unlock()
	return append([]byte(nil), ow.buffer.Bytes()...)
}

func (ow *OutputWriter) String() string {
	ow.lock.Lock()
	defer ow.lock.Unlock()
	return ow.buffer.String()
}

// CloseLines closes the Lines channel. Call when done reading lines.
func (ow *OutputWriter) CloseLines() {
	close(ow.Lines)
}
