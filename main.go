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
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Println("Usage: gomerge [OPTIONS] DST SRC")
		fmt.Println("Options")
		flag.PrintDefaults()
	}

	var (
		writeFile = flag.Bool("w", false, "writes result to destination file")
		rmSrc     = flag.Bool("r", false, "removes source after merge")
	)
	flag.Parse()

	files := flag.Args()

	if len(files) != 2 {
		log.Fatal("missing files, dst src")
	}

	dst := flag.Arg(0)
	src := flag.Arg(1)

	var buf bytes.Buffer
	Merge(&buf, load(dst), load(src))

	if !*writeFile {
		os.Stdout.Write(buf.Bytes())
		os.Exit(0)
	}

	os.WriteFile(dst, buf.Bytes(), 0644)
	if *rmSrc {
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
