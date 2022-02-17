package main

import (
	"go/scanner"
	"go/token"
)

func findImports(src []byte) (importLines []string) {
	s, fset := newScanner(src)

	// skip to imports
	lastEnd := scanTo(token.IMPORT, s, fset)

loop:
	for {
		pos, tok, lit := s.Scan()
		//fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
		switch tok {
		case token.LPAREN:
			lastEnd = fset.Position(pos).Offset + len(lit) + 2
		case token.SEMICOLON:
			lastEnd = fset.Position(pos).Offset + len(lit)

		case token.IDENT: // renamed import
			continue loop
		case token.STRING:
			i := fset.Position(pos).Offset + len(lit)
			line := string(src[lastEnd:i])
			importLines = append(importLines, line)
			lastEnd = i
		case token.RPAREN: // end of imports
			break loop
		case token.EOF: // no imports found
			break loop
		}
	}
	return
}

// scanTo will scan to the given token to find. Returns the position
// after the found literal or -1 if EOF
func scanTo(find token.Token, s *scanner.Scanner, fset *token.FileSet) int {
	var lastEnd int
	for {
		pos, tok, lit := s.Scan()
		lastEnd = fset.Position(pos).Offset + len(lit)
		if tok == find {
			break
		}
		if tok == token.EOF {
			return -1
		}
	}
	return lastEnd
}

func newScanner(src []byte) (*scanner.Scanner, *token.FileSet) {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)
	return &s, fset
}
