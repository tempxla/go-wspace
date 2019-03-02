package main

import (
	"fmt"
)

const (
	TAB   = '\t'
	SPACE = ' '
	LF    = '\n'
)

const (
	psh = iota + 1
	dup
	swp
	pop
	add
	sub
	mul
	div
	mod
	sto
	lod
	mrk
	cll
	jmp
	jze
	jne
	ret
	end
	wtn
	wtc
	rdn
	rdc
)

type imp struct {
	cmd int
	arg int
	lbl string
}

func (imp imp) String() string {
	fmt1 := "[%s]"
	fmt2 := "[%s %d]"
	fmt3 := "[%s %d %s]"
	switch imp.cmd {
	case psh:
		return fmt.Sprintf(fmt2, "psh", imp.arg)
	case dup:
		return fmt.Sprintf(fmt1, "dup")
	case swp:
		return fmt.Sprintf(fmt1, "swp")
	case pop:
		return fmt.Sprintf(fmt1, "pop")
	case add:
		return fmt.Sprintf(fmt1, "add")
	case sub:
		return fmt.Sprintf(fmt1, "sub")
	case mul:
		return fmt.Sprintf(fmt1, "mul")
	case div:
		return fmt.Sprintf(fmt1, "div")
	case mod:
		return fmt.Sprintf(fmt1, "mod")
	case sto:
		return fmt.Sprintf(fmt1, "sto")
	case lod:
		return fmt.Sprintf(fmt1, "lod")
	case mrk:
		return fmt.Sprintf("[%s %s]", "mrk", imp.lbl)
	case cll:
		return fmt.Sprintf(fmt3, "cll", imp.arg, imp.lbl)
	case jmp:
		return fmt.Sprintf(fmt3, "jmp", imp.arg, imp.lbl)
	case jze:
		return fmt.Sprintf(fmt3, "jze", imp.arg, imp.lbl)
	case jne:
		return fmt.Sprintf(fmt3, "jne", imp.arg, imp.lbl)
	case ret:
		return fmt.Sprintf(fmt1, "ret")
	case end:
		return fmt.Sprintf(fmt1, "end")
	case wtn:
		return fmt.Sprintf(fmt1, "wtn")
	case wtc:
		return fmt.Sprintf(fmt1, "wtc")
	case rdn:
		return fmt.Sprintf(fmt1, "rdn")
	case rdc:
		return fmt.Sprintf(fmt1, "rdc")
	default:
		return fmt.Sprintf("undefined")
	}
}
