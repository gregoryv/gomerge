package main

import (
	"testing"
)

func TestScanner_FindImports(t *testing.T) {
	t.Run("multiline imports", func(t *testing.T) {
		s, _ := NewFileScanner("testdata/a.go")
		imports := s.FindImports()

		if imports[0] != "\t\"fmt\"" {
			t.Error(imports)
		}
		if imports[1] != "\t\"strings\"" {
			t.Error(imports)
		}
	})

	t.Run("single import", func(t *testing.T) {
		s, _ := NewFileScanner("testdata/b.go")
		imports := s.FindImports()

		if imports[0] != " \"fmt\"" {
			t.Error(imports)
		}
	})

	t.Run("renamed import", func(t *testing.T) {
		s, _ := NewFileScanner("testdata/c.go")
		imports := s.FindImports()

		if imports[0] != "\tcr \"crypto/rand\"" {
			t.Error(imports)
		}
	})
}
