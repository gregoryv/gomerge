package testdata

import (
	"fmt"
	"strings"
)

// y does stuff
func y() {
	letters := []string{"a", "b"}
	fmt.Println(strings.Join(letters, ","))
}
