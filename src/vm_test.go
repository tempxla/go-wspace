package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestVmStack(t *testing.T) {
	src := `SS(psh) STL(+1)            1
			SS(psh) TTSTL(-101)        -5 1
			SLT(swap)                  1 -5
			SLS(dup)                   1 1 -5
			SLS(dup)                   1 1 1 -5
			SLL(pop)                   1 1 -5
			SLL(pop)                   1 -5
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	eval(cd)
	testEqInt(t, 2, len(stack))
	testEqInt(t, 1, stack[len(stack)-1])
	testEqInt(t, -5, stack[len(stack)-2])
}

func TestVmArith(t *testing.T) {
	src := `SS(psh) STL(+1)            1
			SS(psh) TTSTL(-101)        -5 1
			TSSS(add)                  -4
			SS(psh) TTSTL(-101)        -5 -4
			TSST(sub)                  -9
			SS(psh) TTSTL(-101)        -5 -9
			TSSL(mul)                  45
			SS(psh) TTSTL(-101)        -5 45
			SS(psh) STSL(10)           2 -5 45
			TSTS(div)                  -2 45
			TSTT(mod)                  -1
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	eval(cd)
	testEqInt(t, 1, len(stack))
	testEqInt(t, -1, stack[len(stack)-1])
}

func TestVmHeap(t *testing.T) {
	src := `SS(psh) STTL(+11)          3
			SS(psh) TTSTL(-101)        -5 3
			TTS(sto)                   #{3:-5}
			SS(psh) STSSL(+100)        4
			SS(psh) TTSTSL(-1010)      -10 4
			TTS(sto)                   #{3:-5, 4:-10}
			SS(psh) STTL(+11)          3
			TTT(lod)                   -5
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	eval(cd)
	testEqInt(t, 1, len(stack))
	testEqInt(t, -5, stack[len(stack)-1])
	testEqInt(t, 2, len(heap))
	testEqInt(t, -5, heap[3])
	testEqInt(t, -10, heap[4])
}

func TestVmFlow1(t *testing.T) {
	src := `SS(psh) SSL(+0)         0
			SLS(dup)                0 0
			LTS(jze) SL
			SLL(pop)
			LSS(mrk) SL
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	eval(cd)
	testEqInt(t, 1, len(stack))
}

func TestVmFlow2(t *testing.T) {
	src := `SS(psh) TTL(-1)         -1
			SLS(dup)
			LTT(jne) SL
			SLL(pop)
			LSS(mrk) SL
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	eval(cd)
	testEqInt(t, 1, len(stack))
}

func TestVmFlow3(t *testing.T) {
	src := `SS(psh) STTTL(+7)         7
			SLS(dup)
			LSL(jne) SL
			SLL(pop)
			LSS(mrk) SL
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	eval(cd)
	testEqInt(t, 2, len(stack))
}

func TestVmFlow4(t *testing.T) {
	src := `SS(psh) STSL(+2)         2
			LSL(jmp) TL
			LSS(mrk) SL
			SLS(dup)                 2 2
			LTL(ret)
			LSS(mrk) TL
			LST(cll) SL
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	eval(cd)
	testEqInt(t, 2, len(stack))
}

func TestVmIo1(t *testing.T) {
	src := `SS(psh)   STSSSSSTL(+65)         65(A)
			SLS(dup)
			SLS(dup)
			TLSS(wtc)
			TLST(wtn)
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	buf := &bytes.Buffer{}
	writer = buf
	eval(cd)
	testEqStr(t, "A65", buf.String())
}

func TestVmIo2(t *testing.T) {
	src := `SS(psh)   STSSSSSTL(+65)         65(A)
			TLTS(rdc)
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	buf := &bytes.Buffer{}
	buf.WriteString("B")
	reader = bufio.NewReader(buf)
	eval(cd)
	testEqInt(t, 2, len(stack))
	testEqInt(t, 66, stack[len(stack)-1])
	testEqInt(t, 65, stack[len(stack)-2])
}

func TestVmIo3(t *testing.T) {
	src := `SS(psh)   STSSSSSTL(+65)         65(A)
			TLTT(rdn)
			LLL               `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	buf := &bytes.Buffer{}
	buf.WriteString("15")
	reader = bufio.NewReader(buf)
	eval(cd)
	testEqInt(t, 2, len(stack))
	testEqInt(t, 15, stack[len(stack)-1])
	testEqInt(t, 65, stack[len(stack)-2])
}
