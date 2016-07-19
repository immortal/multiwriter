package main

import (
	"github.com/immortal/multiwriter"
	"os"
)

func main() {
	w1, e := os.Create("file1.txt")
	if e != nil {
		panic(e)
	}

	w2, e := os.Create("file2.txt")
	if e != nil {
		panic(e)
	}

	mw := multiwriter.New(w1, w2)
	data := []byte("Hello ")
	_, e = mw.Write(data)
	if e != nil {
		panic(e)
	}

	var m *multiwriter.MultiWriter = mw.(*multiwriter.MultiWriter)

	m.Remove(w2)
	w2.Close()

	w3, e := os.Create("file3.txt")
	if e != nil {
		panic(e)
	}

	m.Append(w3)

	data = []byte("World ")
	_, e = mw.Write(data)
	if e != nil {
		panic(e)
	}

	w3.Close()
	w1.Close()
}
