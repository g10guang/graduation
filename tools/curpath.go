package tools

import (
	"fmt"
	"path/filepath"
)

func Abs() {
	p, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	fmt.Println(p)
}
