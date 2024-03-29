// Command gomerge merges two or more go files, removing duplicate
// imports.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/gregoryv/gomerge"
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Println("Usage: gomerge [OPTIONS] DST SRC")
		fmt.Println("Options")
		flag.PrintDefaults()
	}

	var (
		writeFile = flag.Bool("w", true, "writes result to destination file")
		rmSrc     = flag.Bool("r", true, "removes source after merge(only with -w)")
		incFile   = flag.Bool("i", false, "include src filename in merged as comment")
	)
	flag.Parse()

	files := flag.Args()

	if len(files) != 2 {
		log.Fatal("missing files, dst src")
	}

	dst := flag.Arg(0)
	src := flag.Arg(1)

	var buf bytes.Buffer
	cmd := gomerge.New(&buf, load(dst), load(src))
	cmd.SetIncludeFile(*incFile)
	cmd.SetSrcFile(src)

	_ = cmd.Run()

	if !*writeFile {
		os.Stdout.Write(buf.Bytes())
		return
	}

	os.WriteFile(dst, buf.Bytes(), 0644)
	if *rmSrc {
		// try git rm -f first
		exec.Command("git", "rm", "-f", src).Run()
		os.RemoveAll(src)
	}
}

func load(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
