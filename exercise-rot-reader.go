package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(b []byte) (int, error) {
	n, err := r.r.Read(b)
	for i:=0;i<n;i++ {
		start := byte(0)
		if b[i] >= 'A' && b[i] <= 'Z' {
			start = 'A'
		} else if b[i] >= 'a' && b[i] <= 'z' {
			start = 'a'
		}
		if start != 0 {
			b[i] = (b[i] - start + 13) % 26 + start
		}
	}

	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
