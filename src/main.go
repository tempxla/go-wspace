package main

import (
	"fmt"
	"io/ioutil"
)

type intStack []int

type state int

const (
	ok = iota
	ng
	end
)

const (
	TAB   = '\t'
	SPACE = ' '
	LF    = '\n'
)

var stack intStack
var heap map[int]int
var callstack intStack
var labels map[int]int
var source []byte
var pointer int

var errMsg string

func init() {
	stack = make([]int, 1024)
	heap = make(map[int]int)
	callstack = make([]int, 1024)
	labels = make(map[int]int)
	pointer = 0
}

func read() byte {
	for { // skip comment
		c := source[pointer]
		pointer++
		switch c {
		case SPACE, TAB, LF:
			return c
		}
	}
}

func impStack() state {
	c := read()
	switch c {
	case SPACE: // [Space] Number Push the number onto the stack
		return number()
	case LF:
		c = read()
		switch c {
		case SPACE: // [LF][Space] Duplicate the top item on the stack
			stack.dup()
		case TAB: //[LF][Tab] Swap the top two items on the stack
			stack.swap()
		case LF: //[LF][LF] Discard the top item on the stack
			stack.discard()
		}
	case TAB:
		return parseError(c)
	}
	return ok
}

func number() state {
	var sign int
	c := read()
	switch c {
	case SPACE: // [Space] for positive
		sign = 1
	case TAB: // [Tab] for negative
		sign = -1
	case LF:
		return parseError(c)
	}
	n, st := label()
	if st == ng {
		return ng
	}
	stack.push(sign * n)
	return ok
}

func (stack *intStack) push(e int) {
	*stack = append(*stack, e)
}

func (stack *intStack) pop() int {
	val := (*stack)[len(*stack)-1]
	(*stack).discard()
	return val
}

func (stack *intStack) dup() {
	*stack = append(*stack, (*stack)[len(*stack)-1])
}

func (stack *intStack) swap() {
	end := len(*stack) - 1
	(*stack)[end-1], (*stack)[end] = (*stack)[end], (*stack)[end-1]
}

func (stack *intStack) discard() {
	*stack = (*stack)[:len(*stack)-2]
}

func (stack *intStack) peek() int {
	return (*stack)[len(*stack)-1]
}

func impArith() state {
	c := read()
	switch c {
	case SPACE:
		c = read()
		switch c {
		case SPACE: // [Space][Space] Addition
			biOp(func(a, b int) int { return a + b })
		case TAB: // [Space][Tab] Subtraction
			biOp(func(a, b int) int { return a - b })
		case LF: // [Space][LF] Multiplication
			biOp(func(a, b int) int { return a * b })
		}
	case TAB:
		c = read()
		switch c {
		case SPACE: // [Tab][Space] Integer Division
			biOp(func(a, b int) int { return a / b })
		case TAB: // [Tab][Tab] Modulo
			biOp(func(a, b int) int { return a % b })
		}
	case LF:
		return parseError(c)
	}
	return ok
}

func biOp(op func(int, int) int) {
	l := stack.pop()
	r := stack.pop()
	stack.push(op(l, r))
}

func impHeap() state {
	c := read()
	switch c {
	case SPACE: // [Space] Store
		val := stack.pop()
		addr := stack.pop()
		heap[addr] = val
	case TAB: // [Tab] Retrieve
		addr := stack.pop()
		stack.push(heap[addr])
	case LF:
		return parseError(c)
	}
	return ok
}

func impFlow() state {
	c := read()
	switch c {
	case SPACE:
		c = read()
		switch c {
		case SPACE: // [Space][Space] Label Mark a location in the program
			mark()
		case TAB: // [Space][Tab] Label Call a subroutine
			call()
		case LF: //[Space][LF] Label Jump unconditionally to a label
			jump()
		}
	case TAB:
		c = read()
		switch c {
		case SPACE: // [Tab][Space] Label Jump to a label if the top of the stack is zero
			jumpZE()
		case TAB: // [Tab][Tab] Label Jump to a label if the top of the stack is negative
			jumpNE()
		case LF: // [Tab][LF] End a subroutine and transfer control back to the caller
			retrun()
		}
	case LF:
		c = read()
		switch c {
		case LF: // [LF][LF] End the program
			return end
		case TAB:
			fallthrough
		case SPACE:
			return parseError(c)
		}
	}
	return ok
}

func label() (int, state) {
	var n int
	c := read()
	switch c {
	case SPACE: // [Space] represents the binary digit 0
		n = 0
	case TAB: // [Tab] represents 1
		n = 1
	case LF:
		return 0, parseError(c)
	}
	for c = read(); c != LF; c = read() {
		switch c {
		case SPACE: // [Space] represents the binary digit 0
			n = n * 2
		case TAB: // [Tab] represents 1
			n = n*2 + 1
		}
	}
	return n, ok
}

func mark() state {
	n, st := label()
	if st == ng {
		return ng
	}
	labels[n] = pointer
	return ok
}

func call() state {
	n, st := label()
	if st == ng {
		return ng
	}
	p := labels[n]
	callstack.push(pointer)
	pointer = p
	return ok
}

func jump() state {
	n, st := label()
	if st == ng {
		return ng
	}
	pointer = labels[n]
	return ok
}

func jumpZE() state {
	n, st := label()
	if st == ng {
		return ng
	}
	if stack.peek() == 0 {
		pointer = labels[n]
	}
	return ok
}

func jumpNE() state {
	n, st := label()
	if st == ng {
		return ng
	}
	if stack.peek() < 0 {
		pointer = labels[n]
	}
	return ok
}

func retrun() {
	pointer = callstack.pop()
}

func impIO() state {
	c := read()
	switch c {
	case SPACE:
		c = read()
		switch c {
		case SPACE: // [Space][Space] Output the character at the top of the stack
			// TODO
		case TAB: // [Space][Tab] Output the number at the top of the stack
			// TODO
		case LF:
			return parseError(c)
		}
	case TAB:
		c = read()
		switch c {
		case SPACE: // [Tab][Space] Read a character and place it in the location given by the top of the stack
			// TODO
		case TAB: // [Tab][Tab] Read a number and place it in the location given by the top of the stack
			// TODO
		case LF:
			return parseError(c)
		}
	case LF:
		return parseError(c)
	}
	return ok
}

func eval() state {
	var st state
	c := read()
	switch c {
	case SPACE: // [Space] Stack Manipulation
		st = impStack()
	case TAB:
		c = read()
		switch c {
		case SPACE: //[Tab][Space] Arithmetic
			st = impArith()
		case TAB: // [Tab][Tab] Heap access
			st = impHeap()
		case LF: // [Tab][LF] I/O
			st = impIO()
		default:
			st = parseError(c)
		}
	case LF: // [LF] Flow Control
		st = impFlow()
	}
	switch st {
	case ok:
		eval()
	case ng:
		return ng
	default:
		return end
	}
	return end
}

func parseError(c byte) state {
	errMsg = fmt.Sprintf("error: unexpected character [%c], pos:%d\n", c, pointer)
	return ng
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
