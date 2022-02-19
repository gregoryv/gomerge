package main

import (
	"bytes"
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
