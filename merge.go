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

	// todo assumes d has imports
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
	if s.HasImports() {
		MoveAfterImports(s)
	} else {
		s.ScanTo(token.PACKAGE)
		s.ScanTo(token.STRING)
	}
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

func FindImports(s *Scanner) (importLines []string) {
	// skip to imports
	s.ScanTo(token.IMPORT)

loop:
	for {
		pos, tok, lit := s.Scan()
		//fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
		switch tok {
		case token.LPAREN:
			s.lastEnd = s.fset.Position(pos).Offset + len(lit) + 2
		case token.SEMICOLON:
			s.lastEnd = s.fset.Position(pos).Offset + len(lit)

		case token.IDENT: // renamed import
			continue loop
		case token.STRING:
			i := s.fset.Position(pos).Offset + len(lit)
			line := string(s.src[s.lastEnd:i])
			importLines = append(importLines, line)
			s.lastEnd = i
		case token.RPAREN: // end of imports
			break loop
		case token.EOF: // no imports found
			break loop
		}
	}
	return
}

func MoveAfterImports(s *Scanner) {
	if !s.HasImports() {
		return
	}
	s.ScanTo(token.IMPORT)
	_, tok, _ := s.Scan()
	switch tok {
	case token.LPAREN:
		s.ScanTo(token.RPAREN)
	case token.STRING:
		s.ScanTo(token.SEMICOLON)
	}
}
