package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	dir := t.TempDir()
	data, _ := os.ReadFile(c)
	x := filepath.Join(dir, "x.go")
	os.WriteFile(x, data, 0644)

	os.Args = []string{"gomerge", "-w", a, b, x}
	main()

	out, _ := os.ReadFile(x)
	checkMerge(t, out)
}

func Test_Merge(t *testing.T) {

	var buf bytes.Buffer
	Merge(&buf, files)
	checkMerge(t, buf.Bytes())
}

func checkMerge(t *testing.T, result []byte) {
	t.Helper()

	got := string(result)
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

const (
	a = "testdata/a.go"
	b = "testdata/b.go"
	c = "testdata/c.go"
)

var (
	files = []string{a, b, c}
)
