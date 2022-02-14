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

	dstFile := files[len(files)-1]
	dst, err := os.ReadFile(dstFile)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	for _, srcFile := range files[:len(files)-1] {
		src, err := os.ReadFile(srcFile)
		if err != nil {
			log.Fatal(err)
		}
		if err := MergeGoFiles(&buf, dst, src); err != nil {
			log.Fatal(err)
		}
	}
	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if !writeToFile {
		os.Stdout.Write(out)
		os.Exit(0)
	}
	os.WriteFile(dstFile, out, 0644)
}

// MergeGoFiles merges file a into file b
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
					destImports.Specs = append(destImports.Specs, iSpec)
				}
			}
		}
	}
}
