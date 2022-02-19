package main

import (
	"go/token"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	t.Run("plain", func(t *testing.T) {
		data := []byte(`package x`)
		s := Split(data)
		if s.Header.Len() == 0 {
			t.Error("empty header")
		}
		if s.Imports.Len() != 0 {
			t.Error("found imports: ", s.Imports.String())
		}
		if s.Rest.Len() != 0 {
			t.Error("found rest: ", s.Rest.String())
		}
	})

	t.Run("one import", func(t *testing.T) {
		data := []byte(`package x

import "fmt"
`)
		s := Split(data)
		if s.Header.Len() == 0 {
			t.Error("empty header")
		}
		if strings.Contains(s.Header.String(), "import") {
			t.Error("found import in header")
		}
		if s.Imports.Len() == 0 {
			t.Error("empty imports")
		}
		if s.Rest.Len() != 0 {
			t.Error("found rest: ", s.Rest.String())
		}
	})

	t.Run("one import and body", func(t *testing.T) {
		data := []byte(`package x

import "fmt"

func x() {}
`)
		s := Split(data)
		if s.Header.Len() == 0 {
			t.Error("empty header")
		}
		if strings.Contains(s.Header.String(), "import") {
			t.Error("found import in header")
		}
		if s.Imports.Len() == 0 {
			t.Error("empty imports")
		}
		if s.Rest.Len() == 0 {
			t.Error("empty rest")
		}
	})
}

func TestScanner_FindImports(t *testing.T) {
	t.Run("multiline imports", func(t *testing.T) {
		s, _ := NewFileScanner("testdata/a.go")
		imports := FindImports(s)

		if imports[0] != "\t\"fmt\"" {
			t.Error(imports)
		}
		if imports[1] != "\t\"strings\"" {
			t.Error(imports)
		}
	})

	t.Run("single import", func(t *testing.T) {
		s, _ := NewFileScanner("testdata/b.go")
		imports := FindImports(s)

		if imports[0] != " \"fmt\"" {
			t.Error(imports)
		}
	})

	t.Run("renamed import", func(t *testing.T) {
		s, _ := NewFileScanner("testdata/c.go")
		imports := FindImports(s)

		if imports[0] != "\tcr \"crypto/rand\"" {
			t.Error(imports)
		}
	})
}

func TestScanner_ScanTo(t *testing.T) {
	s, _ := NewFileScanner("testdata/a.go")
	block, _ := s.ScanTo(token.IMPORT)
	got := string(block)
	if !strings.Contains(got, "package testdata") {
		t.Error(got)
	}
	if !strings.Contains(got, "import") {
		t.Error(got)
	}
}
