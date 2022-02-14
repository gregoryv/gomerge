package main

import (
	"bytes"
	"go/format"
	"go/token"
	"io"
	"log"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
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
	a, err := decorator.Parse(fileA)
	if err != nil {
		return err
	}

	b, err := decorator.Parse(fileB)
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

	// write b including new imports
	if err := decorator.Fprint(w, b); err != nil {
		return err
	}

	return nil
}

func merge(dest, src *dst.File) {
	// find dst import declaration
	var destImports *dst.GenDecl
	for i := 0; i < len(dest.Decls); i++ {
		d := dest.Decls[i]

		switch d.(type) {
		case *dst.FuncDecl:
			// No action
		case *dst.GenDecl:
			dd := d.(*dst.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				destImports = dd
			}
		}
	}

	for i := 0; i < len(src.Decls); i++ {
		d := src.Decls[i]

		switch d.(type) {
		case *dst.FuncDecl:
			dest.Decls = append(dest.Decls, d)

		case *dst.GenDecl:
			dd := d.(*dst.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				// skip
				for _, iSpec := range src.Imports {
					destImports.Specs = append(destImports.Specs, iSpec)
				}
			}
		}
	}
}
