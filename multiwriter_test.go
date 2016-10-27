package multiwriter

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type brokenWriter struct{}
type halfBrokenWriter struct{}

func (bw *brokenWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("This is supposed to break")
}

func (bw *halfBrokenWriter) Write(p []byte) (n int, err error) {
	return 2, nil
}

func TestNew(t *testing.T) {
	out := new(bytes.Buffer)
	mw := New(out)
	multi, ok := mw.(*MultiWriter)
	if !ok {
		t.Fatal("Could not convert writer to Multi Writer")
	}
	if len(multi.writers) != 1 {
		t.Fatal("Not the same number of writers")
	}
}

func TestWriter(t *testing.T) {
	out := new(bytes.Buffer)
	out2 := new(bytes.Buffer)
	mw := New(out, out2)
	msg := []byte("This is a test")
	n, err := mw.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg) != n {
		t.Fatal("Did not write the ammount of bytes needed")
	}
	if string(msg) != out.String() || string(msg) != out2.String() {
		t.Fatal("Did not write the correct message.")
	}
}

func TestWriter_Broken(t *testing.T) {
	out := new(bytes.Buffer)
	out2 := new(brokenWriter)
	mw := New(out, out2)
	msg := []byte("This is a test")
	n, err := mw.Write(msg)
	if err == nil {
		t.Fatal("I expect an error from a broken writer")
	}
	if n != 0 {
		t.Fatal("The broken writer should not have written any bytes")
	}
	if err.Error() != "This is supposed to break" {
		t.Fatal("Not the expected error from writer")
	}
}

func TestWriter_HalfBroken(t *testing.T) {
	out := new(bytes.Buffer)
	out2 := new(halfBrokenWriter)
	mw := New(out, out2)
	msg := []byte("This is a test")
	n, err := mw.Write(msg)
	if err == nil {
		t.Fatal("I expect an error from a half broken writer")
	}
	if n != 2 {
		t.Fatal("The broken writer should have written only two bits")
	}
	if err.Error() != io.ErrShortWrite.Error() {
		t.Fatal("Not the expected error from writer")
	}
}

func TestLen(t *testing.T) {
	out := new(bytes.Buffer)
	out2 := new(bytes.Buffer)
	mw := New(out, out2)
	multi, ok := mw.(*MultiWriter)
	if !ok {
		t.Fatal("Could not convert writer to Multi Writer")
	}
	if multi.Len() != 2 || len(multi.writers) != multi.Len() {
		t.Fatal("Not the same number of writers")
	}
}

func TestAppend(t *testing.T) {
	out := new(bytes.Buffer)
	out2 := new(bytes.Buffer)
	out3 := new(bytes.Buffer)
	msg := []byte("This is a test")
	mw := New(out, out2)
	multi, ok := mw.(*MultiWriter)
	if !ok {
		t.Fatal("Could not convert writer to Multi Writer")
	}
	multi.Append(out3)
	if multi.Len() != 3 {
		t.Fatal("Not the expected number of writers")
	}
	n, err := mw.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg) != n {
		t.Fatal("Did not write the ammount of bytes needed")
	}
	if string(msg) != out.String() || string(msg) != out2.String() || string(msg) != out3.String() {
		t.Fatal("Did not write the correct message.")
	}
}

func TestRemove(t *testing.T) {
	out := new(bytes.Buffer)
	out2 := new(bytes.Buffer)
	out3 := new(bytes.Buffer)
	msg := []byte("This is a test")
	mw := New(out, out2, out3)
	multi, ok := mw.(*MultiWriter)
	if !ok {
		t.Fatal("Could not convert writer to Multi Writer")
	}
	multi.Remove(out2)
	if multi.Len() != 2 {
		t.Fatal("Not the expected number of writers")
	}
	n, err := mw.Write(msg)
	if err != nil {
		t.Fatal(err)
	}
	if len(msg) != n {
		t.Fatal("Did not write the ammount of bytes needed")
	}
	if string(msg) != out.String() || string(msg) != out3.String() {
		t.Fatal("Did not write the correct message.")
	}
	if out2.String() != "" {
		t.Fatal("We still wrote to a removed writer")
	}
}
