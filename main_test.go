package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func Test_MergeGoFiles(t *testing.T) {
	a, _ := os.ReadFile("testdata/a.go")
	b, _ := os.ReadFile("testdata/b.go")
	var buf bytes.Buffer
	if err := MergeGoFiles(&buf, a, b); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	exp := []string{
		"func x()",
		"func y()",
		`"fmt"`,
		`"strings"`,
	}
	for _, exp := range exp {
		if !strings.Contains(got, exp) {
			t.Log(got)
			t.Fatal("missing", exp)
		}
	}
}
