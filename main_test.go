package main

import (
	"bytes"
	"os"
	"os/exec"
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

func Test_BeforeAfter(t *testing.T) {
	before := run(t, "go", "run", "./testdata/before")
	var buf bytes.Buffer
	files, _ := filepath.Glob("./testdata/before/*.go")
	Merge(&buf, files)

	filename := filepath.Join(t.TempDir(), "main.go")
	os.WriteFile(filename, buf.Bytes(), 0644)
	after := run(t, "go", "run", filename)

	if string(before) != string(after) {
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(string(data))
		t.Error("after not the same as before")
	}
}

func run(t *testing.T, app string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command(app, args...)
	out, err := cmd.CombinedOutput()
	t.Log(cmd.String(), "\n", string(out))
	if err != nil {
		t.Error(err)
	}
	return out
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
