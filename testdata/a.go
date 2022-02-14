package testdata

import (
	"fmt"
	"strings"
)

func y() {
	letters := []string{"a", "b"}
	fmt.Println(strings.Join(letters, ","))
}
