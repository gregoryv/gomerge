package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/gregoryv/golden"
)

func TestMerge(t *testing.T) {
	var (
		buf bytes.Buffer
		dst = []byte(`// my pkg
package x
import "fmt"
func main() {}
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
	if err := Merge(&buf, dst, src); err != nil {
		t.Fatal(err)
	}

	golden.Assert(t, buf.String())
}

func ExampleSplit_header() {
	data := []byte(`// my docs
package x
`)
	s := Split(data)
	fmt.Println(s.Header)
	// output:
	// // my docs
}

func ExampleSplit_single() {
	data := []byte(`package x

import "fmt"

func x() {}
`)
	s := Split(data)
	fmt.Println(s.Imports)
	// output:
	// import "fmt"
}

func ExampleSplit() {
	data := []byte(`package x

import (
	"fmt"
)

func x() {}
`)
	s := Split(data)
	fmt.Println(s.Imports)
	// output:
	//	import (
	//	"fmt"
	//)
}

func TestSplit(t *testing.T) {
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
