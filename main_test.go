package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_MergeGoFiles(t *testing.T) {
	files := []string{
		"testdata/a.go",
		"testdata/b.go",
		"testdata/c.go",
	}
	var buf bytes.Buffer
	Merge(&buf, files)

	got := buf.String()
	exp := []string{
		"func x()",
		"func y()",
		`"fmt"`,
		`"strings"`,
		"// y does stuff",
	}
	for _, exp := range exp {
		if !strings.Contains(got, exp) {
			t.Log(got)
			t.Fatal("missing", exp)
		}
	}
	if strings.Count(got, `"fmt"`) > 1 {
		t.Error("multiple fmt imports")
	}
}
