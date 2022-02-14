package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
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
		MergeGoFiles(&buf, dst, src)
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(out)
}

// MergeGoFiles merges fiel a into file b
func MergeGoFiles(w io.Writer, fileB, fileA []byte) error {
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
	merge(b, a)

	ast.SortImports(bfset, b)

	// write b including new imports
	if err := format.Node(w, bfset, b); err != nil {
		return err
	}

	return err
}

func merge(dst, src *ast.File) {
	// find dst import declaration
	var dstImports *ast.GenDecl
	for i := 0; i < len(dst.Decls); i++ {
		d := dst.Decls[i]

		switch d.(type) {
		case *ast.FuncDecl:
			// No action
		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				dstImports = dd
			}
		}
	}

	for i := 0; i < len(src.Decls); i++ {
		d := src.Decls[i]

		switch d.(type) {
		case *ast.FuncDecl:
			dst.Decls = append(dst.Decls, d)

		case *ast.GenDecl:
			dd := d.(*ast.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				// skip
				for _, iSpec := range src.Imports {
					dstImports.Specs = append(dstImports.Specs, iSpec)
				}
			}
		}
	}
}
