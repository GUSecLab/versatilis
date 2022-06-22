package versatilis

import (
	"errors"
	sync "sync"
)

type VBufferOverflowError struct{}

func (e *VBufferOverflowError) Error() string {
	return "overflow occurred"
}

type vBuffer struct {
	init      bool `default:"false"`
	buf       []byte
	mu        *sync.Mutex
	sizeLimit int // or 0 or -1 if no limit
}

func NewVBuffer() *vBuffer {
	buf := &vBuffer{
		init:      true,
		mu:        &sync.Mutex{},
		sizeLimit: -1,
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

// attempts to read the entirety of the buffer.  Returns the bytes read
func (buf *vBuffer) ReadAll() (b []byte, err error) {
	if !buf.init {
		return nil, errors.New("buffer not initialized")
	}
	buf.mu.Lock()
	defer buf.mu.Unlock()

	b = buf.buf
	buf.buf = nil
	return b, nil
}

// writes to the buffer.  Returns the number of bytes written or err on error
func (buf *vBuffer) Write(b []byte) (n int, err error) {
	if !buf.init {
		return -1, errors.New("buffer not initialized")
	}
	buf.mu.Lock()
	defer buf.mu.Unlock()

	if buf.sizeLimit > 0 {
		currentBufSize := len(buf.buf)
		bytesToWrite := buf.sizeLimit - currentBufSize
		switch {
		case bytesToWrite < 0:
			panic("buffer overflow; programming error :(")
		case bytesToWrite == 0:
			return 0, &VBufferOverflowError{}
		default:
			bs := b[:bytesToWrite]
			buf.buf = append(buf.buf, bs...)
			return bytesToWrite, nil
		}
	} else {
		buf.buf = append(buf.buf, b...)
		return len(b), nil
	}
}
