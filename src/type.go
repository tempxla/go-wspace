package main

const (
	TAB   = '\t'
	SPACE = ' '
	LF    = '\n'
)

const (
	psh = iota
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
}
