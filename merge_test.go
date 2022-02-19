package gomerge

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func Example() {
	Merge(os.Stdout,
		[]byte(`package x

import "fmt"

func x() { fmt.Println("hello") }`),
		[]byte(`package x

import "strings"

func y() { strings.Repeat(" ", 10) }`),
	)
	// output:
	// package x
	//
	// import (
	// 	"fmt"
	// 	"strings"
	// )
	//
	// func x() { fmt.Println("hello") }
	//
	// func y() { strings.Repeat(" ", 10) }
}

func TestSplit(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		data := []byte{}
		Split(data)
	})

	t.Run("plain", func(t *testing.T) {
		data := []byte(`package x
func x() {}
`)
		s := Split(data)
		if len(s.Header) != 0 {
			t.Error("found header: ", s.Header)
		}
		if len(s.Package) == 0 {
			t.Error("empty package")
		}
		if len(s.Imports) != 0 {
			t.Error("found imports: ", s.Imports)
		}
		if len(s.Rest) == 0 {
			t.Error("empty rest")
		}
		if s.String() != string(data) {
			t.Error("not equal")
		}
	})

	t.Run("one import", func(t *testing.T) {
		data := []byte(`package x

import "fmt"
`)
		s := Split(data)
		if len(s.Header) != 0 {
			t.Error("found header:", s.Header)
		}
		if strings.Contains(s.Header, "import") {
			t.Error("found import in header")
		}
		if len(s.Imports) == 0 {
			t.Error("empty imports")
		}
		if len(s.Rest) != 0 {
			t.Error("found rest: ", s.Rest)
		}
	})

	t.Run("one import and body", func(t *testing.T) {
		data := []byte(`// my docs
package x

import (
	"fmt"
	"strings"
)

func x() {}
`)
		s := Split(data)
		if len(s.Header) == 0 {
			t.Error("empty header")
		}
		if strings.Contains(s.Header, "import") {
			t.Error("found import in header")
		}
		if len(s.Imports) == 0 {
			t.Error("empty imports")
		}
		if len(s.Rest) == 0 {
			t.Error("empty rest")
		}
	})
}

func TestGoMerge_Run(t *testing.T) {
	var (
		buf bytes.Buffer
		dst = []byte(`// my pkg
package x

import "fmt"
func gomerge() {}
`)
		src = []byte(`// other
package x

import (
"fmt"
"strings"
)

func x() {}
`)
	)

	cmd := GoMerge{
		w:   &buf,
		dst: dst,

		includeFile: true,
		srcFile:     "test",
		src:         src,
	}

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	golden.Assert(t, buf.String())
}
