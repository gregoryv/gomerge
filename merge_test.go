package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/gregoryv/golden"
)

func TestMerge(t *testing.T) {
	var (
		buf    bytes.Buffer
		dst, _ = os.ReadFile("testdata/a.go")
		src, _ = os.ReadFile("testdata/b.go")
	)
	if err := Merge(&buf, dst, src); err != nil {
		t.Fatal(err)
	}
	t.Log("dst:", string(dst))
	t.Log("\nsrc:", string(src))
	t.Log("----------")
	golden.Assert(t, buf.String())
}
