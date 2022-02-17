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

// ScanTo will scan to the given token to find. Returns the content
// since last scan and position after the found literal or -1 if EOF
// The result includes the token to find.
func (s *Scanner) ScanTo(find token.Token) ([]byte, int) {
	start := s.lastEnd
	for {
		pos, tok, lit := s.Scan()
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

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	pos, tok, lit = s.Scanner.Scan()
	s.current = pos
	return
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
