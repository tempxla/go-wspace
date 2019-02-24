package main

import (
	"fmt"
	"io/ioutil"
)

const cTAB = '\t'
const cSPACE = ' '
const cLF = '\n'

var stack []int
var heap map[int]int
var callstack []int
var source []byte
var pointer int

func init() {
	stack = make([]int, 1024)
	heap = make(map[int]int)
	callstack = make([]int, 1024)
	pointer = 0
}

func read() byte {
	c := source[pointer]
	pointer++
	return c
}

func impStack() {
	c := read()
	switch c {
	case cSPACE: // [Space] NumberPush the number onto the stack
		number()
	case cLF:
		c = read()
		switch c {
		case cSPACE: // [LF][Space] Duplicate the top item on the stack
			dup()
		case cTAB: //[LF][Tab] Swap the top two items on the stack
			swap()
		case cLF: //[LF][LF] Discard the top item on the stack
			discard()
		default:
			parseError(c)
		}
	default:
		parseError(c)
	}
}

func number() {
	var sign int
	num := 0
	c := read()
	switch c {
	case cSPACE: // [Space] for positive
		sign = 1
	case cTAB: // [Tab] for negative
		sign = -1
	default:
		parseError(c)
	}
	for c := read(); c != cLF; {
		switch c {
		case cSPACE: // [Space] represents the binary digit 0
			num = num * 2
		case cTAB: // [Tab] for negative
			num = num*2 + 1
		default:
			parseError(c)
		}
	}
	push(sign * num)
}

func push(e int) {
	stack = append(stack, e)
}

func pop() int {
	val := stack[len(stack)-1]
	discard()
	return val
}

func dup() {
	stack = append(stack, stack[len(stack)-1])
}

func swap() {
	end := len(stack) - 1
	stack[end-1], stack[end] = stack[end], stack[end-1]
}

func discard() {
	stack = stack[:len(stack)-2]
}

func impArith() {
	c := read()
	switch c {
	case cSPACE:
		c = read()
		switch c {
		case cSPACE: // [Space][Space] Addition
			biOp(func(a, b int) int { return a + b })
		case cTAB: // [Space][Tab] Subtraction
			biOp(func(a, b int) int { return a - b })
		case cLF: // [Space][LF] Multiplication
			biOp(func(a, b int) int { return a * b })
		default:
			parseError(c)
		}
	case cTAB:
		c = read()
		switch c {
		case cSPACE: // [Tab][Space] Integer Division
			biOp(func(a, b int) int { return a / b })
		case cTAB: // [Tab][Tab] Modulo
			biOp(func(a, b int) int { return a % b })
		default:
			parseError(c)
		}
	default:
		parseError(c)
	}
}

func biOp(op func(int, int) int) {
	l := pop()
	r := pop()
	push(op(l, r))
}

func impHeap() {
	c := read()
	switch c {
	case cSPACE: // [Space] Store
		val := pop()
		addr := pop()
		heap[addr] = val
	case cTAB: // [Tab] Retrieve
		addr := pop()
		push(heap[addr])
	default:
		parseError(c)
	}
}

func impFlow() {
	c := read()
	switch c {
	case cSPACE:
		c = read()
		switch c {
		case cSPACE: // [Space][Space] Label Mark a location in the program
			// TODO
		case cTAB: // [Space][Tab] Label Call a subroutine
			// TODO
		case cLF: //[Space][LF] Label Jump unconditionally to a label
			// TODO
		default:
			parseError(c)
		}
	case cTAB:
		c = read()
		switch c {
		case cSPACE: // [Tab][Space] Label Jump to a label if the top of the stack is zero
			// TODO
		case cTAB: // [Tab][Tab] Label Jump to a label if the top of the stack is negative
			// TODO
		case cLF: // [Tab][LF] End a subroutine and transfer control back to the caller
			// TODO
		default:
			parseError(c)
		}
	case cLF:
		c = read()
		switch c {
		case cLF: // [LF][LF] End the program
			// TODO
		default:
			parseError(c)
		}
	default:
		parseError(c)
	}
}

func impIO() {
	c := read()
	switch c {
	case cSPACE:
		c = read()
		switch c {
		case cSPACE: // [Space][Space] Output the character at the top of the stack
			// TODO
		case cTAB: // [Space][Tab] Output the number at the top of the stack
			// TODO
		default:
			parseError(c)
		}
	case cTAB:
		c = read()
		switch c {
		case cSPACE: // [Tab][Space] Read a character and place it in the location given by the top of the stack
			// TODO
		case cTAB: // [Tab][Tab] Read a number and place it in the location given by the top of the stack
			// TODO
		default:
			parseError(c)
		}
	default:
		parseError(c)
	}
}

func eval() {
	c := read()
	switch c {
	case cSPACE: // [Space] Stack Manipulation
		impStack()
	case cTAB:
		c = read()
		switch c {
		case cSPACE: //[Tab][Space] Arithmetic
			impArith()
		case cTAB: // [Tab][Tab] Heap access
			impHeap()
		case cLF: // [Tab][LF] I/O
			impIO()
		default:
			parseError(c)
		}
	case cLF: // [LF] Flow Control
		impFlow()
	default:
		parseError(c)
	}
}

func parseError(c byte) {
	fmt.Printf("error: unexpected character [%c]\n", c)
}

func parseFromFile(path string) {
	var err error
	source, err = ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	eval()
}

func main() {
	parseFromFile("")
}
