package multiwriter

import "io"

type multiWriter struct {
	writers []io.Writer
}

// Write implements io.Writer
func (self *multiWriter) Write(p []byte) (n int, err error) {
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

// New return a multiWriter
func New(writers ...io.Writer) io.Writer {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return &multiWriter{w}
}
