package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	x, _, _ := reader.ReadRune()
	fmt.Println(x)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	y := scanner.Text()
	fmt.Println(y)
}
