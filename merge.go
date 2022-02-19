package main

import (
	"bytes"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"strings"
)

func Merge(w io.Writer, dst, src []byte) error {
	d := Split(dst)
	s := Split(src)

	d.Header.WriteTo(w)
	d.Package.WriteTo(w)
	fmt.Fprint(w, "\n")

	imports := mergeImports(
		d.Imports.Bytes(),
		s.Imports.Bytes(),
	)
	w.Write(imports)
	fmt.Fprint(w, "\n")

	d.Rest.WriteTo(w)
	s.Rest.WriteTo(w)
	return nil
}

func mergeImports(a, b []byte) []byte {
	var buf bytes.Buffer
	all := append(
		importLines(a),
		importLines(b)...,
	)
	// todo filter duplicates
	buf.WriteString("import (\n")
	for _, line := range unique(all) {
		buf.WriteString("\t")
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	buf.WriteString(")")
	return buf.Bytes()
}

func importLines(expr []byte) []string {
	lines := strings.Split(string(expr), "\n")

	res := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
		case line == ")":
		case strings.HasPrefix(line, "import ("):
		case strings.HasPrefix(line, "import "):
			res = append(res, line[7:])
		default:
			res = append(res, line)
		}
	}
	return res
}

func unique(v []string) []string {
	h := make(map[string]int)
	res := make([]string, 0)
	for _, v := range v {
		if _, found := h[v]; found {
			continue
		}
		h[v] = 1
		res = append(res, v)
	}
	return res
}

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
