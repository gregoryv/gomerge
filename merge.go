package main

import "io"

func Merge(w io.Writer, dst, src []byte) error {
	d := NewScanner(dst)
	s := NewScanner(src)

	_ = d
	_ = s
	return nil
}
