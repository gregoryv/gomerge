// Command gomerge merges two or more go files, removing duplicate
// imports.
package main

import (
	"bytes"
	"flag"
	"fmt"
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

	flag.Usage = func() {
		fmt.Println("Usage: gomerge [OPTION] SRC... DST")
		flag.PrintDefaults()
	}

	var writeToFile bool
	flag.BoolVar(&writeToFile, "w", writeToFile, "writes result to destination file")
	flag.Parse()

	files := flag.Args()
	if len(files) < 2 {
		log.Fatal("missing files, ...src dst")
	}

	var buf bytes.Buffer
	Merge(&buf, files)

	if !writeToFile {
		os.Stdout.Write(buf.Bytes())
		os.Exit(0)
	}
	dstFile := files[len(files)-1]
	os.WriteFile(dstFile, buf.Bytes(), 0644)
}

func Merge(w io.Writer, files []string) {
	// load destination
	dstFile := files[len(files)-1]
	dstData, err := os.ReadFile(dstFile)
	shouldNot(err)

	dest, err := decorator.Parse(dstData)
	shouldNot(err)

	for _, srcFile := range files[:len(files)-1] {
		srcData, err := os.ReadFile(srcFile)
		shouldNot(err)
		src, err := decorator.Parse(srcData)
		shouldNot(err)

		merge(dest, src)
	}

	var buf bytes.Buffer
	// write out destination source
	err = decorator.Fprint(&buf, dest)
	shouldNot(err)

	// tidy, todo use gofmt or not at all
	out, err := format.Source(buf.Bytes())
	shouldNot(err)
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

	exists := func(s *dst.ImportSpec) bool {
		for _, d := range destImports.Specs {
			d := d.(*dst.ImportSpec)
			if s.Path.Value == d.Path.Value {
				return true
			}
		}
		return false
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
					if !exists(iSpec) {
						destImports.Specs = append(destImports.Specs, iSpec)
					}
				}
			}
		}
	}
}

func shouldNot(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}
