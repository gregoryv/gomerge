package main

import (
	"bytes"
	"strings"
	"testing"
)

func Test_Merge(t *testing.T) {
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
		"func z()",
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

	// single occurence
	unique := []string{
		`"fmt"`,
		"package testdata",
	}
	for _, s := range unique {
		if strings.Count(got, s) > 1 {
			t.Error("multiple:", s)
		}
	}
}
