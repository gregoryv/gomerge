package main

import (
	"fmt"
	"go/token"
	"io"
	"strings"
)

func Merge(w io.Writer, dst, src []byte) error {
	d := NewScanner(dst)
	s := NewScanner(src)

	d.ScanTo(token.IMPORT)
	w.Write(d.Before())

	imports := mergeStrings(
		FindImports(NewScanner(dst)),
		FindImports(NewScanner(src)),
	)

	fmt.Fprint(w, "import (\n")
	for _, imp := range imports {
		fmt.Fprintf(w, "\t%s\n", imp)
	}
	fmt.Fprint(w, ")")

	// TODO find end of optional imports
	d.ScanTo(token.RPAREN)
	w.Write(d.Rest())

	// puts it at end of imports
	s.ScanTo(token.IMPORT)
	s.ScanTo(token.RPAREN)

	w.Write(s.Rest())

	return nil
}

func mergeStrings(a, b []string) []string {
	existing := make(map[string]int)
	unique := make([]string, 0)
	for _, line := range append(a, b...) {
		line = strings.TrimSpace(line)
		if _, found := existing[line]; found {
			continue
		}
		existing[line] = 1
		unique = append(unique, line)
	}
	return unique
}
