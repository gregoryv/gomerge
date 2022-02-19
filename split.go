package main

import (
	"bytes"
	"go/scanner"
	"go/token"
)

func Split(src []byte) *GoSrc {
	gos := GoSrc{}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	var (
		buf          = &gos.Header
		i            int
		packageFound bool
	)
loop:
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		j := fset.Position(pos).Offset
		if tok == token.PACKAGE {
			buf.Write(src[i:j])
			i = j
			buf = &gos.Package
		}

		switch tok {
		case token.STRING, token.IDENT:
			j = fset.Position(pos).Offset + len(lit)
		default:
			j = fset.Position(pos).Offset + len(tok.String())
		}

		//log.Println("i:", i, "j:", j)
		if j >= len(src) {
			buf.Write(src[i:])
		} else {
			buf.Write(src[i:j])
		}
		i = j
		switch {
		case tok == token.SEMICOLON && packageFound:
			break loop
		case tok == token.PACKAGE:
			packageFound = true
		}
	}

	// scan for imports
	buf = &gos.Imports
	var endFound, multi bool
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		var j int
		switch tok {
		case token.STRING, token.IDENT:
			j = fset.Position(pos).Offset + len(lit)
		case token.LPAREN:
			multi = true
			j = fset.Position(pos).Offset + len(tok.String())
		default:
			j = fset.Position(pos).Offset + len(tok.String())
		}

		//log.Println("i:", i, "j:", j)
		if j >= len(src) {
			buf.Write(src[i:])
		} else {
			buf.Write(src[i:j])
		}
		i = j
		if endFound {
			break
		}
		if multi && tok == token.RPAREN {
			endFound = true
			continue
		}
		if !multi && tok == token.SEMICOLON {
			break
		}
	}

	// and the rest
	buf = &gos.Rest
	if i < len(src) {
		buf.Write(src[i:])
	}
	return &gos
}

type GoSrc struct {
	Header  bytes.Buffer // docs before package
	Package bytes.Buffer // package
	Imports bytes.Buffer // imports
	Rest    bytes.Buffer // rest of the content
}

func (s *GoSrc) String() string {
	var buf bytes.Buffer
	s.Header.WriteTo(&buf)
	s.Package.WriteTo(&buf)
	s.Imports.WriteTo(&buf)
	s.Rest.WriteTo(&buf)
	return buf.String()
}
