package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

var code []imp
var addr int
var stack intStack
var heap map[int]int
var callstack intStack
var reader *bufio.Reader
var writer io.Writer

func eval(cd []imp) {
	code = cd
	addr = 0
	stack = make([]int, 0, 1024)
	heap = make(map[int]int)
	callstack = make([]int, 0, 1024)
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	if writer == nil {
		writer = os.Stdout
	}
	// for i, imp := range code {
	// 	fmt.Printf("%4d %v\n", i, imp)
	// }
	runEval()
}

func runEval() {
	//fmt.Println(":::addr::: ", addr, stack)
	cd := code[addr]
	addr++
	switch cd.cmd {
	case psh:
		stack.push(cd.arg)
	case dup:
		a := stack.pop()
		stack.push(a)
		stack.push(a)
	case swp:
		a := stack.pop()
		b := stack.pop()
		stack.push(a)
		stack.push(b)
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
		// do nothing
	case cll:
		callstack.push(addr)
		addr = cd.arg
	case jmp:
		addr = cd.arg
	case jze:
		if stack.pop() == 0 {
			addr = cd.arg
		}
	case jne:
		if stack.pop() < 0 {
			addr = cd.arg
		}
	case ret:
		addr = callstack.pop()
	case end:
		return
	case wtn:
		fmt.Fprint(writer, stack.pop())
	case wtc:
		fmt.Fprint(writer, string(stack.pop()))
	case rdn:
		s, _ := reader.ReadString('\n')
		i, err := strconv.Atoi(strings.TrimRight(s, "\r\n"))
		if err != nil {
			fmt.Println(err)
		} else {
			heap[stack.pop()] = i
		}
	case rdc:
		r, _, _ := reader.ReadRune()
		heap[stack.pop()] = int(r)
	}
	runEval()
}
