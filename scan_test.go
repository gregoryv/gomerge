package main

import (
	"os"
	"testing"
)

func Test_findImports(t *testing.T) {
	t.Run("multiline imports", func(t *testing.T) {
		src, _ := os.ReadFile("testdata/a.go")
		imports := findImports(src)

		if imports[0] != "\t\"fmt\"" {
			t.Error(imports)
		}
		if imports[1] != "\t\"strings\"" {
			t.Error(imports)
		}
	})

	t.Run("single import", func(t *testing.T) {
		src, _ := os.ReadFile("testdata/b.go")
		imports := findImports(src)

		if imports[0] != " \"fmt\"" {
			t.Error(imports)
		}
	})

	t.Run("renamed import", func(t *testing.T) {
		src, _ := os.ReadFile("testdata/c.go")
		imports := findImports(src)

		if imports[0] != "\tcr \"crypto/rand\"" {
			t.Error(imports)
		}
	})
}
