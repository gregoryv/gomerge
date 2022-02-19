package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func Merge(w io.Writer, dst, src []byte) error {
	d := Split(dst)
	s := Split(src)

	d.Header.WriteTo(w)
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
	for _, line := range all {
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
