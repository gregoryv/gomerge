package testdata

import (
	"fmt"
	"strings"
)

type Car struct{}

// y does stuff
func y() {
	letters := []string{"a", "b"}
	fmt.Println(strings.Join(letters, ","))
}
