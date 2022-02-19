package gomerge

import (
	"bytes"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"strings"
)

func Merge(w io.Writer, dst, src []byte) error {
	cmd := &GoMerge{
		w:   w,
		dst: dst,
		src: src,
	}
	return cmd.Run()
}

func New(w io.Writer, dst, src []byte) *GoMerge {
	return &GoMerge{
		w:   w,
		dst: dst,
		src: src,
	}
}

type GoMerge struct {
	w   io.Writer
	dst []byte

	includeFile bool
	srcFile     string
	src         []byte
}

func (me *GoMerge) Run() error {
	d := Split(me.dst)
	s := Split(me.src)
	w := me.w

	fmt.Fprint(w, d.Header)
	fmt.Fprint(w, d.Package)
	fmt.Fprint(w, "\n")

	imports := mergeImports(
		[]byte(d.Imports),
		[]byte(s.Imports),
	)
	w.Write(imports)
	fmt.Fprint(w, "\n")

	fmt.Fprint(w, d.Rest)

	if me.includeFile {
		fmt.Fprintln(w, "\n// gomerge src:", me.srcFile)
	}
	fmt.Fprint(w, s.Header)
	fmt.Fprint(w, s.Rest)
	return nil
}

func (me *GoMerge) SetIncludeFile(v bool) {
	me.includeFile = v
}

func (me *GoMerge) SetSrcFile(v string) {
	me.srcFile = v
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
		buf          bytes.Buffer
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
			gos.Header = buf.String()
			buf.Reset()
			buf.Write(src[i:j])
			i = j
		}

		switch tok {
		case token.STRING, token.IDENT, token.COMMENT:
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
	gos.Package = buf.String()
	buf.Reset()

	// scan for imports
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
	gos.Imports = buf.String()
	buf.Reset()

	// and the rest
	if i < len(src) {
		buf.Write(src[i:])
	}

	gos.Rest = buf.String()
	return &gos
}

type GoSrc struct {
	Header  string // docs before package
	Package string // package
	Imports string // imports
	Rest    string // rest of the content
}

func (s *GoSrc) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.Header)
	buf.WriteString(s.Package)
	buf.WriteString(s.Imports)
	buf.WriteString(s.Rest)
	return buf.String()
}
