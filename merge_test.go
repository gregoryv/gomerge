package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
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
	got := buf.String()
	if got == "" {
		t.Error("empty")
	}

	if strings.Count(got, "package") == 2 {
		t.Error("duplicate package")
	}
}
