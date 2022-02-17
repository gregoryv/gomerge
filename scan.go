package main

import (
	"bytes"
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
	current token.Pos
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

// ScanTo will scan to the given token to find. Returns the content
// since last scan and position after the found literal or -1 if EOF
// The result includes the token to find.
func (s *Scanner) ScanTo(find token.Token) ([]byte, int) {
	start := s.lastEnd
	for {
		pos, tok, lit := s.Scan()
		s.current = pos
		switch tok {
		case token.STRING, token.IDENT:
			s.lastEnd = s.fset.Position(pos).Offset + len(lit)
		default:
			s.lastEnd = s.fset.Position(pos).Offset + len(tok.String())
		}
		if tok == find {
			break
		}
		if tok == token.EOF {
			return s.src[start:], -1
		}
	}
	return s.src[start:s.lastEnd], s.lastEnd
}

func (s *Scanner) Rest() []byte {
	return s.src[s.lastEnd:]
}

func (s *Scanner) Before() []byte {
	return s.src[:s.fset.Position(s.current).Offset]
}

func (s *Scanner) HasImports() bool {
	return bytes.Index(s.src, []byte("import")) > 0
}
