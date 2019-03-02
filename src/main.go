package main

import (
	"fmt"
	"os"
)

func main() {
	cd, err := parseFromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	eval(cd)
}
