package main

import (
	"go/scanner"
	"go/token"
	"os"
)

func NewFileScanner(filename string) (*Scanner, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewScanner(src), nil
}

func NewScanner(src []byte) *Scanner {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)
	return &Scanner{
		Scanner: &s,
		fset:    fset,
		src:     src,
	}
}

type Scanner struct {
	*scanner.Scanner
	fset    *token.FileSet
	src     []byte
	lastEnd int
}

func (s *Scanner) FindImports() (importLines []string) {
	// skip to imports
	s.scanTo(token.IMPORT)

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

// scanTo will scan to the given token to find. Returns the position
// after the found literal or -1 if EOF
func (s *Scanner) scanTo(find token.Token) int {
	for {
		pos, tok, lit := s.Scan()
		s.lastEnd = s.fset.Position(pos).Offset + len(lit)
		if tok == find {
			break
		}
		if tok == token.EOF {
			return -1
		}
	}
	return s.lastEnd
}
