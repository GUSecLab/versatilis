package versatilis

import (
	"errors"
	sync "sync"
)

type vBuffer struct {
	init bool `default:"false"`
	buf  []byte
	mu   *sync.Mutex
}

func NewBuffer() *vBuffer {
	buf := &vBuffer{
		init: true,
		mu:   &sync.Mutex{},
	}
	return buf
}

// returns the size of the buffer
func (buf *vBuffer) Size() int {
	if buf == nil || !buf.init {
		return 0
	} else {
		buf.mu.Lock()
		defer buf.mu.Unlock()
		return len(buf.buf)
	}
}

// attempts to read numBytes from the buffer.  Returns the bytes read
func (buf *vBuffer) Read(numBytes int) (b []byte, err error) {
	if !buf.init {
		return nil, errors.New("buffer not initialized")
	}
	buf.mu.Lock()
	defer buf.mu.Unlock()

	n := minInt(len(buf.buf), numBytes)
	if n == 0 {
		return nil, nil
	}
	b = make([]byte, n)
	b = buf.buf[:n]
	buf.buf = buf.buf[n:]
	return b, nil
}

// writes to the buffer.  Returns the number of bytes written or err on error
func (buf *vBuffer) Write(b []byte) (n int, err error) {
	if !buf.init {
		return -1, errors.New("buffer not initialized")
	}
	buf.mu.Lock()
	defer buf.mu.Unlock()

	buf.buf = append(buf.buf, b...)
	return len(b), nil
}
