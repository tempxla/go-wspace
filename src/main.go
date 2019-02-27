package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type intStack []int

func (stack *intStack) push(e int) {
	*stack = append(*stack, e)
}

func (stack *intStack) pop() int {
	val := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return val
}

func (stack *intStack) swap() {
	end := len(*stack) - 1
	(*stack)[end-1], (*stack)[end] = (*stack)[end], (*stack)[end-1]
}

func (stack *intStack) peek() int {
	return (*stack)[len(*stack)-1]
}

var stack intStack
var heap map[int]int
var callstack intStack
var code []imp
var addr int
var reader *bufio.Reader

func eval() {
	cd := code[addr]
	addr++
	switch cd.cmd {
	case psh:
		stack.push(cd.arg)
	case dup:
		stack.push(stack.peek())
	case swp:
		stack.swap()
	case pop:
		stack.pop()
	case add:
		r, l := stack.pop(), stack.pop()
		stack.push(l + r)
	case sub:
		r, l := stack.pop(), stack.pop()
		stack.push(l - r)
	case mul:
		r, l := stack.pop(), stack.pop()
		stack.push(l * r)
	case div:
		r, l := stack.pop(), stack.pop()
		stack.push(l / r)
	case mod:
		r, l := stack.pop(), stack.pop()
		stack.push(l % r)
	case sto:
		v, addr := stack.pop(), stack.pop()
		heap[addr] = v
	case lod:
		stack.push(heap[stack.pop()])
	case mrk:
		// unreachable
	case cll:
		callstack.push(addr)
		addr = cd.arg
	case jmp:
		addr = cd.arg
	case jze:
		if stack.peek() == 0 {
			addr = cd.arg
		}
	case jne:
		if stack.peek() < 0 {
			addr = cd.arg
		}
	case ret:
		addr = callstack.pop()
	case end:
		return
	case wtn:
		fmt.Print(stack.peek())
	case wtc:
		fmt.Print(rune(stack.peek()))
	case rdn:
		s, _ := reader.ReadString('\n')
		i, _ := strconv.Atoi(s)
		stack.push(i)
	case rdc:
		r, _, _ := reader.ReadRune()
		stack.push(int(r))
	}
	eval()
}

func main() {
	var err error
	code, err = parseFromFile("")
	if err != nil {
		fmt.Println(err)
		return
	}
	addr = 0
	stack = make([]int, 1024)
	heap = make(map[int]int)
	callstack = make([]int, 1024)
	reader = bufio.NewReader(os.Stdin)
	eval()
}
