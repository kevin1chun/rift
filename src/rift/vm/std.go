package vm

import (
	"fmt"
	"strings"
)

func println(objs...interface{}) {
	var strs []string
	for _, obj := range objs {
		strs = append(strs, fmt.Sprintf("%s", obj))
	}
	fmt.Println(strings.Join(strs, " "))
}