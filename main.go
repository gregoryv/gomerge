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
		fmt.Println("Usage: gomerge [OPTION] SRC... DST")
		flag.PrintDefaults()
	}

	var writeToFile bool
	flag.BoolVar(&writeToFile, "w", writeToFile, "writes result to destination file")
	flag.Parse()

	files := flag.Args()

	if len(files) < 2 {
		log.Fatal("missing files, ...src dst")
	}

	var buf bytes.Buffer
	//Merge(&buf, files)

	if !writeToFile {
		os.Stdout.Write(buf.Bytes())
		os.Exit(0)
	}
	dstFile := files[len(files)-1]
	os.WriteFile(dstFile, buf.Bytes(), 0644)
}
