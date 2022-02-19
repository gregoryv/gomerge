package gomerge_test

import (
	"os"

	"github.com/gregoryv/gomerge"
)

func Example() {
	cmd := gomerge.New(os.Stdout,
		[]byte(`package x

import "fmt"

func x() { fmt.Println("hello") }`),
		[]byte(`package x

import "strings"

func y() { strings.Repeat(" ", 10) }`),
	)
	cmd.SetSrcFile("mytest.go")
	cmd.SetIncludeFile(true)
	cmd.Run()
	// output:
	// package x
	//
	// import (
	// 	"fmt"
	// 	"strings"
	// )
	//
	// func x() { fmt.Println("hello") }
	// // gomerge src: mytest.go
	//
	// func y() { strings.Repeat(" ", 10) }
}
