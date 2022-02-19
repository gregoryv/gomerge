package main

import (
	"fmt"
	"strings"
	"testing"
)

func ExampleSplit_single() {
	data := []byte(`package x

import "fmt"

func x() {}
`)
	s := Split(data)
	fmt.Println(s.Imports.String())
	// output:
	// "fmt"
}

func ExampleSplit() {
	data := []byte(`package x

import (
	"fmt"
)

func x() {}
`)
	s := Split(data)
	fmt.Println(s.Imports.String())
	// output:
	//	import (
	//	"fmt"
	//)
}

func TestSplit(t *testing.T) {
	t.Run("plain", func(t *testing.T) {
		data := []byte(`package x`)
		s := Split(data)
		if s.Header.Len() == 0 {
			t.Error("empty header")
		}
		if s.Imports.Len() != 0 {
			t.Error("found imports: ", s.Imports.String())
		}
		if s.Rest.Len() != 0 {
			t.Error("found rest: ", s.Rest.String())
		}
	})

	t.Run("one import", func(t *testing.T) {
		data := []byte(`package x

import "fmt"
`)
		s := Split(data)
		if s.Header.Len() == 0 {
			t.Error("empty header")
		}
		if strings.Contains(s.Header.String(), "import") {
			t.Error("found import in header")
		}
		if s.Imports.Len() == 0 {
			t.Error("empty imports")
		}
		if s.Rest.Len() != 0 {
			t.Error("found rest: ", s.Rest.String())
		}
	})

	t.Run("one import and body", func(t *testing.T) {
		data := []byte(`package x

import (
	"fmt"
	"strings"
)

func x() {}
`)
		s := Split(data)
		if s.Header.Len() == 0 {
			t.Error("empty header")
		}
		if strings.Contains(s.Header.String(), "import") {
			t.Error("found import in header")
		}
		if s.Imports.Len() == 0 {
			t.Error("empty imports")
		}
		if s.Rest.Len() == 0 {
			t.Error("empty rest")
		}
	})
}
