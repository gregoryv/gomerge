// Command gomerge merges two or more go files, removing duplicate
// imports.
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

func Merge(w io.Writer, files []string) {
	// load destination
	dest := load(files[len(files)-1])

	for _, srcFile := range files[:len(files)-1] {
		src := load(srcFile)
		merge(dest, src)
	}

	var buf bytes.Buffer
	// write out destination source
	check(decorator.Fprint(&buf, dest))

	// tidy, todo use gofmt or not at all
	out, err := format.Source(buf.Bytes())
	check(err)
	w.Write(out)
}

func merge(dest, src *dst.File) {
	// find dst import declaration
	var destImports *dst.GenDecl
	for i := 0; i < len(dest.Decls); i++ {
		d := dest.Decls[i]

		switch d.(type) {
		case *dst.GenDecl:
			dd := d.(*dst.GenDecl)
			if dd.Tok == token.IMPORT {
				destImports = dd
				break
			}
		}
	}

	// copy declarations
	for i := 0; i < len(src.Decls); i++ {
		d := src.Decls[i]

		switch d.(type) {
		case *dst.FuncDecl:
			dest.Decls = append(dest.Decls, d)

		case *dst.GenDecl:
			dd := d.(*dst.GenDecl)

			// IMPORT Declarations are grouped
			if dd.Tok == token.IMPORT {
				// skip
				for _, iSpec := range src.Imports {
					if !exists(destImports, iSpec) {
						destImports.Specs = append(destImports.Specs, iSpec)
					}
				}
				return
			}
			dest.Decls = append(dest.Decls, d)
		}
	}
}

// returns true if import exists in destination import declaration
func exists(dest *dst.GenDecl, s *dst.ImportSpec) bool {
	for _, d := range dest.Specs {
		d := d.(*dst.ImportSpec)
		if s.Path.Value == d.Path.Value {
			return true
		}
	}
	return false
}

func load(filename string) *dst.File {
	// load destination
	data, err := os.ReadFile(filename)
	check(err)

	defer func() {
		e := recover()
		if e != nil {
			log.Fatal("invalid go file: ", filename)
		}
	}()
	f, err := decorator.Parse(data)
	check(err)
	return f
}

func check(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}
