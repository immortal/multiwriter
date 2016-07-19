package multiwriter

import (
	"io"
	"sync"
)

type MultiWriter struct {
	sync.Mutex
	writers []io.Writer
}

// Append writer
func (self *MultiWriter) Append(writers ...io.Writer) {
	self.Lock()
	defer self.Unlock()
	self.writers = append(self.writers, writers...)
}

// Remove writer
func (self *MultiWriter) Remove(writers ...io.Writer) {
	self.Lock()
	defer self.Unlock()
	for i := len(self.writers) - 1; i > 0; i-- {
		for _, v := range writers {
			if self.writers[i] == v {
				self.writers = append(self.writers[:i], self.writers[i+1:]...)
				break
			}
		}
	}
}

// Write implements io.Writer
func (self *MultiWriter) Write(p []byte) (n int, err error) {
	self.Lock()
	defer self.Unlock()

	type result struct {
		n   int
		err error
	}

	rs := make(chan *result)

	for _, w := range self.writers {
		go func(writer io.Writer) {
			n, err := writer.Write(p)
			rs <- &result{n, err}
		}(w)
	}

	for range self.writers {
		r := <-rs
		if r.err != nil {
			return r.n, r.err
		}
		if r.n != len(p) {
			return r.n, io.ErrShortWrite
		}
	}
	return len(p), nil
}

// New return a MultiWriter
func New(writers ...io.Writer) io.Writer {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return &MultiWriter{writers: w}
}
