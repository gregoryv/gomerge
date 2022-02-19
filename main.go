// Command gomerge merges two or more go files, removing duplicate
// imports.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Println("Usage: gomerge [OPTION] DST SRC")
		flag.PrintDefaults()
	}

	var writeToFile bool
	flag.BoolVar(&writeToFile, "w", writeToFile, "writes result to destination file")
	flag.Parse()

	files := flag.Args()

	if len(files) != 2 {
		log.Fatal("missing files, dst src")
	}

	dst := flag.Arg(0)
	src := flag.Arg(1)

	var buf bytes.Buffer
	Merge(&buf, load(dst), load(src))

	if !writeToFile {
		os.Stdout.Write(buf.Bytes())
		os.Exit(0)
	}

	os.WriteFile(dst, buf.Bytes(), 0644)
}

func load(filename string) []byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
