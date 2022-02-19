package gomerge

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func TestSplit(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		data := []byte{}
		Split(data)
	})

	t.Run("plain", func(t *testing.T) {
		data := []byte(`package x`)
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
		if len(s.Rest) != 0 {
			t.Error("found rest: ", s.Rest)
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

func TestMerge(t *testing.T) {
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
