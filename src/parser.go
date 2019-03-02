package main

import (
	"fmt"
	"io/ioutil"
)

type state struct {
	src   []byte
	pos   int
	addrs map[string]int // label -> address
	code  []imp
}

type parseError struct {
	act byte
	pos int
	msg string
}

func newParseError(act byte, pos int) *parseError {
	return &parseError{act: act, pos: pos}
}

func (e *parseError) Error() string {
	if e.msg == "" {
		return fmt.Sprintf("error: unexpected character [%d], pos:%d\n", e.act, e.pos)
	} else {
		return e.msg
	}
}

func read(st *state) (c byte, eof bool) {
	if st.pos == len(st.src) {
		eof = true
		return
	}
	for st.pos < len(st.src) { // skip comments
		c = st.src[st.pos]
		st.pos++
		switch c {
		case SPACE, TAB, LF:
			return c, false
		}
	}
	return
}

func impStack(st *state) (err error) {
	c, eof := read(st)
	if eof {
		return
	}
	switch c {
	case SPACE: // [Space] Number Push the number onto the stack
		return number(st)
	case LF:
		c, eof = read(st)
		if eof {
			return
		}
		switch c {
		case SPACE: // [LF][Space] Duplicate the top item on the stack
			st.code = append(st.code, imp{cmd: dup})
		case TAB: //[LF][Tab] Swap the top two items on the stack
			st.code = append(st.code, imp{cmd: swp})
		case LF: //[LF][LF] Discard the top item on the stack
			st.code = append(st.code, imp{cmd: pop})
		}
	case TAB:
		return newParseError(c, st.pos)
	}
	return
}

func number(st *state) (err error) {
	var sign int
	c, eof := read(st)
	if eof {
		return
	}
	switch c {
	case SPACE: // [Space] for positive
		sign = 1
	case TAB: // [Tab] for negative
		sign = -1
	case LF:
		err = newParseError(c, st.pos)
		return
	}
	n, err := uint(st)
	if err != nil {
		return
	}
	st.code = append(st.code, imp{cmd: psh, arg: sign * n})
	return
}

func uint(st *state) (n int, err error) {
	c, eof := read(st)
	if eof {
		return
	}
	n = 0
	switch c {
	case SPACE: // [Space] represents the binary digit 0
		n = n * 2
	case TAB: // [Tab] represents 1
		n = n*2 + 1
	case LF:
		err = newParseError(c, st.pos)
		return
	}
	for c, eof = read(st); c != LF && !eof; c, eof = read(st) {
		switch c {
		case SPACE: // [Space] represents the binary digit 0
			n = n * 2
		case TAB: // [Tab] represents 1
			n = n*2 + 1
		}
	}
	return
}

func impArith(st *state) (err error) {
	c, eof := read(st)
	if eof {
		return
	}
	switch c {
	case SPACE:
		c, eof = read(st)
		if eof {
			return
		}
		switch c {
		case SPACE: // [Space][Space] Addition
			st.code = append(st.code, imp{cmd: add})
		case TAB: // [Space][Tab] Subtraction
			st.code = append(st.code, imp{cmd: sub})
		case LF: // [Space][LF] Multiplication
			st.code = append(st.code, imp{cmd: mul})
		}
	case TAB:
		c, eof = read(st)
		if eof {
			return
		}
		switch c {
		case SPACE: // [Tab][Space] Integer Division
			st.code = append(st.code, imp{cmd: div})
		case TAB: // [Tab][Tab] Modulo
			st.code = append(st.code, imp{cmd: mod})
		}
	case LF:
		return newParseError(c, st.pos)
	}
	return
}

func impHeap(st *state) (err error) {
	c, eof := read(st)
	if eof {
		return
	}
	switch c {
	case SPACE: // [Space] Store
		st.code = append(st.code, imp{cmd: sto})
	case TAB: // [Tab] Retrieve
		st.code = append(st.code, imp{cmd: lod})
	case LF:
		return newParseError(c, st.pos)
	}
	return
}

func impFlow(st *state) (err error) {
	c, eof := read(st)
	if eof {
		return
	}
	switch c {
	case SPACE:
		c, eof = read(st)
		if eof {
			return
		}
		l, err := label(st)
		if err != nil {
			return err
		}
		switch c {
		case SPACE: // [Space][Space] Label Mark a location in the coderam
			if _, exist := st.addrs[l]; exist {
				return &parseError{msg: fmt.Sprintf("Label Duplicated: %s", l)}
			}
			st.code = append(st.code, imp{cmd: mrk, lbl: l})
			st.addrs[l] = len(st.code)
		case TAB: // [Space][Tab] Label Call a subroutine
			st.code = append(st.code, imp{cmd: cll, lbl: l})
		case LF: //[Space][LF] Label Jump unconditionally to a label
			st.code = append(st.code, imp{cmd: jmp, lbl: l})
		}
	case TAB:
		c, eof = read(st)
		if eof {
			return
		}
		if c == LF { // [Tab][LF] End a subroutine and transfer control back to the caller
			st.code = append(st.code, imp{cmd: ret})
			return
		}
		l, err := label(st)
		if err != nil {
			return err
		}
		switch c {
		case SPACE: // [Tab][Space] Label Jump to a label if the top of the stack is zero
			st.code = append(st.code, imp{cmd: jze, lbl: l})
		case TAB: // [Tab][Tab] Label Jump to a label if the top of the stack is negative
			st.code = append(st.code, imp{cmd: jne, lbl: l})
		}
	case LF:
		c, eof = read(st)
		if eof {
			return
		}
		switch c {
		case LF: // [LF][LF] End the coderam
			st.code = append(st.code, imp{cmd: end})
		case TAB, SPACE:
			return newParseError(c, st.pos)
		}
	}
	return
}

func label(st *state) (label string, err error) {
	c, eof := read(st)
	if eof {
		return
	}
	var bs []byte
	switch c {
	case SPACE: // [Space] represents the binary digit 0
		bs = append(bs, 'S')
	case TAB: // [Tab] represents 1
		bs = append(bs, 'T')
	case LF:
		err = newParseError(c, st.pos)
		return
	}
	for c, eof = read(st); c != LF && !eof; c, eof = read(st) {
		switch c {
		case SPACE: // [Space] represents the binary digit 0
			bs = append(bs, 'S')
		case TAB: // [Tab] represents 1
			bs = append(bs, 'T')
		}
	}
	label = string(bs)
	return
}

func impIO(st *state) (err error) {
	c, eof := read(st)
	if eof {
		return
	}
	switch c {
	case SPACE:
		c, eof = read(st)
		if eof {
			return
		}
		switch c {
		case SPACE: // [Space][Space] Output the character at the top of the stack
			st.code = append(st.code, imp{cmd: wtc})
		case TAB: // [Space][Tab] Output the number at the top of the stack
			st.code = append(st.code, imp{cmd: wtn})
		case LF:
			return newParseError(c, st.pos)
		}
	case TAB:
		c, eof = read(st)
		if eof {
			return
		}
		switch c {
		case SPACE: // [Tab][Space] Read a character and place it in the location given by the top of the stack
			st.code = append(st.code, imp{cmd: rdc})
		case TAB: // [Tab][Tab] Read a number and place it in the location given by the top of the stack
			st.code = append(st.code, imp{cmd: rdn})
		case LF:
			return newParseError(c, st.pos)
		}
	case LF:
		return newParseError(c, st.pos)
	}
	return
}

func runParser(st *state) (err error) {
	c, eof := read(st)
	if eof {
		return
	}
	switch c {
	case SPACE: // [Space] Stack Manipulation
		err = impStack(st)
	case TAB:
		c, eof = read(st)
		if eof {
			return
		}
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
		return
	}
	return runParser(st)
}

func parseFromFile(path string) (code []imp, err error) {
	src, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return parse(src)
}

func parse(src []byte) (code []imp, err error) {
	st := state{
		src:   src,
		pos:   0,
		addrs: make(map[string]int),
		code:  make([]imp, 0, 1024),
	}
	err = runParser(&st)
	if err != nil {
		return
	}
	// label -> address
	for i := range st.code {
		switch st.code[i].cmd {
		case cll, jmp, jze, jne:
			addr, ok := st.addrs[st.code[i].lbl]
			if !ok {
				err = &parseError{msg: fmt.Sprintf("Label Not Found: %s", st.code[i].lbl)}
				return
			}
			st.code[i].arg = addr
		}
	}
	code = st.code
	return
}
