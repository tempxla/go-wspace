package main

import (
	"fmt"
	"io/ioutil"
)

type state struct {
	src   []byte
	pos   int
	addrs map[int]int // label -> address
	code  []imp
}

type parseError struct {
	act byte
	pos int
}

func (e *parseError) Error() string {
	return fmt.Sprintf("error: unexpected character [%c], pos:%d\n", e.act, e.pos)
}

func read(st *state) byte {
	for { // skip comments
		c := st.src[st.pos]
		(*st).pos++
		switch c {
		case SPACE, TAB, LF:
			return c
		}
	}
}

func impStack(st *state) (err error) {
	c := read(st)
	switch c {
	case SPACE: // [Space] Number Push the number onto the stack
		return number(st)
	case LF:
		c = read(st)
		switch c {
		case SPACE: // [LF][Space] Duplicate the top item on the stack
			st.code = append(st.code, imp{cmd: dup})
		case TAB: //[LF][Tab] Swap the top two items on the stack
			st.code = append(st.code, imp{cmd: swp})
		case LF: //[LF][LF] Discard the top item on the stack
			st.code = append(st.code, imp{cmd: pop})
		}
	case TAB:
		return &parseError{}
	}
	return
}

func number(st *state) (err error) {
	var sign int
	c := read(st)
	switch c {
	case SPACE: // [Space] for positive
		sign = 1
	case TAB: // [Tab] for negative
		sign = -1
	case LF:
		err = &parseError{}
		return
	}
	n, err := label(st)
	if err != nil {
		return
	}
	st.code = append(st.code, imp{cmd: psh, arg: sign * n})
	return
}

func impArith(st *state) (err error) {
	c := read(st)
	switch c {
	case SPACE:
		c = read(st)
		switch c {
		case SPACE: // [Space][Space] Addition
			st.code = append(st.code, imp{cmd: add})
		case TAB: // [Space][Tab] Subtraction
			st.code = append(st.code, imp{cmd: sub})
		case LF: // [Space][LF] Multiplication
			st.code = append(st.code, imp{cmd: mul})
		}
	case TAB:
		c = read(st)
		switch c {
		case SPACE: // [Tab][Space] Integer Division
			st.code = append(st.code, imp{cmd: div})
		case TAB: // [Tab][Tab] Modulo
			st.code = append(st.code, imp{cmd: mod})
		}
	case LF:
		return &parseError{}
	}
	return
}

func impHeap(st *state) (err error) {
	c := read(st)
	switch c {
	case SPACE: // [Space] Store
		st.code = append(st.code, imp{cmd: sto})
	case TAB: // [Tab] Retrieve
		st.code = append(st.code, imp{cmd: lod})
	case LF:
		return &parseError{}
	}
	return
}

func impFlow(st *state) (err error) {
	c := read(st)
	switch c {
	case SPACE:
		c = read(st)
		n, err := label(st)
		if err != nil {
			return err
		}
		switch c {
		case SPACE: // [Space][Space] Label Mark a location in the coderam
			st.addrs[n] = len(st.code)
		case TAB: // [Space][Tab] Label Call a subroutine
			st.code = append(st.code, imp{cmd: cll, arg: n})
		case LF: //[Space][LF] Label Jump unconditionally to a label
			st.code = append(st.code, imp{cmd: jmp, arg: n})
		}
	case TAB:
		c = read(st)
		n, err := label(st)
		if err != nil {
			return err
		}
		switch c {
		case SPACE: // [Tab][Space] Label Jump to a label if the top of the stack is zero
			st.code = append(st.code, imp{cmd: jze, arg: n})
		case TAB: // [Tab][Tab] Label Jump to a label if the top of the stack is negative
			st.code = append(st.code, imp{cmd: jne, arg: n})
		case LF: // [Tab][LF] End a subroutine and transfer control back to the caller
			st.code = append(st.code, imp{cmd: ret})
		}
	case LF:
		c = read(st)
		switch c {
		case LF: // [LF][LF] End the coderam
			st.code = append(st.code, imp{cmd: end})
		case TAB:
			fallthrough
		case SPACE:
			return &parseError{}
		}
	}
	return nil
}

func label(st *state) (n int, err error) {
	c := read(st)
	switch c {
	case SPACE: // [Space] represents the binary digit 0
		n = 0
	case TAB: // [Tab] represents 1
		n = 1
	case LF:
		err = &parseError{}
		return
	}
	for c = read(st); c != LF; c = read(st) {
		switch c {
		case SPACE: // [Space] represents the binary digit 0
			n = n * 2
		case TAB: // [Tab] represents 1
			n = n*2 + 1
		}
	}
	return
}

func impIO(st *state) (err error) {
	c := read(st)
	switch c {
	case SPACE:
		c = read(st)
		switch c {
		case SPACE: // [Space][Space] Output the character at the top of the stack
			st.code = append(st.code, imp{cmd: wtc})
		case TAB: // [Space][Tab] Output the number at the top of the stack
			st.code = append(st.code, imp{cmd: wtn})
		case LF:
			return &parseError{}
		}
	case TAB:
		c = read(st)
		switch c {
		case SPACE: // [Tab][Space] Read a character and place it in the location given by the top of the stack
			st.code = append(st.code, imp{cmd: rdc})
		case TAB: // [Tab][Tab] Read a number and place it in the location given by the top of the stack
			st.code = append(st.code, imp{cmd: rdn})
		case LF:
			return &parseError{}
		}
	case LF:
		return &parseError{}
	}
	return
}

func parse(st *state) (err error) {
	c := read(st)
	switch c {
	case SPACE: // [Space] Stack Manipulation
		err = impStack(st)
	case TAB:
		c = read(st)
		switch c {
		case SPACE: // [Tab][Space] Arithmetic
			err = impArith(st)
		case TAB: // [Tab][Tab] Heap access
			err = impHeap(st)
		case LF: // [Tab][LF] I/O
			err = impIO(st)
		}
	case LF: // [LF] Flow Control
		err = impFlow(st)
	}
	if err != nil {
		return err
	}
	parse(st)
	return
}

func parseFromFile(path string) (code []imp, err error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	st := state{src: source, pos: 0, addrs: make(map[int]int), code: make([]imp, 1024)}
	err = parse(&st)
	if err != nil {
		return
	}
	// label -> address
	for _, cd := range st.code {
		switch cd.cmd {
		case cll, jmp, jze, jne:
			cd.arg = st.addrs[cd.arg]
		}
	}
	return
}
