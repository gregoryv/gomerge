package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	files := os.Args[1:]
	if len(files) < 2 {
		log.Fatal("missing files, ...src dst")
	}

	dstFile := files[len(files)-1]
	dst, _ := os.ReadFile(dstFile)

	var buf bytes.Buffer
	for _, srcFile := range files[:len(files)-1] {
		src, _ := os.ReadFile(srcFile)
		MergeGoFiles(&buf, src, dst)
	}
	io.Copy(os.Stdout, &buf)
	fmt.Fprintln(os.Stdout)
}

// MergeGoFiles merges fiel a into file b
func MergeGoFiles(w io.Writer, fileA, fileB []byte) error {
	afset := token.NewFileSet()
	a, err := parser.ParseFile(afset, "", fileA, 0)
	if err != nil {
		return err
	}

	bfset := token.NewFileSet()
	b, err := parser.ParseFile(bfset, "", fileB, 0)
	if err != nil {
		return err
	}

	// Start writing output
	imports := make(map[string]struct{})

	// find distinct imports
	for _, s := range b.Imports {
		imports[s.Path.Value] = struct{}{}
	}

	// copy missing to fileB, won't work
	var s *ast.ImportSpec
	for _, s = range a.Imports {
		if _, found := imports[s.Path.Value]; !found {
			addImport(b, s)
		}
	}
	lastImport := s

	ast.SortImports(bfset, b)

	// write b including new imports
	if err := printer.Fprint(w, bfset, b); err != nil {
		return err
	}

	var buf bytes.Buffer
	format.Node(&buf, afset, a)
	restOfA := buf.Bytes()[lastImport.End()+1:]
	w.Write([]byte("\n"))
	_, err = w.Write(bytes.TrimSpace(restOfA))
	return err
}

func addImport(dst *ast.File, iSpec *ast.ImportSpec) {
	for i := 0; i < len(dst.Decls); i++ {
		d := dst.Decls[i]

		switch d.(type) {
		case *ast.FuncDecl:
			// No action
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				dd.Specs = append(dd.Specs, iSpec)
			}
		}
	}
}
