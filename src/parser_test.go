package main

import (
	"testing"
)

const errorFormat = "expected:%v, actual:%v"

func testEqInt(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf(errorFormat, expected, actual)
	}
}

func testEqImp(t *testing.T, expected imp, actual imp) {
	if expected != actual {
		t.Errorf(errorFormat, expected, actual)
	}
}

func mkSource(str string) []byte {
	var bs []byte
	for _, c := range str {
		switch c {
		case 'S':
			bs = append(bs, SPACE)
		case 'T':
			bs = append(bs, TAB)
		case 'L':
			bs = append(bs, LF)
		}
	}
	return bs
}

func TestParseStack(t *testing.T) {
	src := `SS(push) STL(+1)
			SS(push) TTSTL(-101)
			SLS(dup)
			SLT(swp)
			SLL(pop)
			LLL`
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	testEqInt(t, 6, len(cd))
	testEqImp(t, imp{cmd: psh, arg: 1}, cd[0])
	testEqImp(t, imp{cmd: psh, arg: -5}, cd[1])
	testEqImp(t, imp{cmd: dup}, cd[2])
	testEqImp(t, imp{cmd: swp}, cd[3])
	testEqImp(t, imp{cmd: pop}, cd[4])
	testEqImp(t, imp{cmd: end}, cd[5])
}

func TestParseArith(t *testing.T) {
	src := `TSSS(add)
			TSST(sub)
			TSSL(mul)
			TSTS(div)
			TSTT(mod)
			LLL`
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	testEqInt(t, 6, len(cd))
	testEqImp(t, imp{cmd: add}, cd[0])
	testEqImp(t, imp{cmd: sub}, cd[1])
	testEqImp(t, imp{cmd: mul}, cd[2])
	testEqImp(t, imp{cmd: div}, cd[3])
	testEqImp(t, imp{cmd: mod}, cd[4])
	testEqImp(t, imp{cmd: end}, cd[5])
}

func TestParseHeap(t *testing.T) {
	src := `TTS(sto)
			TTT(lod)
			LLL`
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	testEqInt(t, 3, len(cd))
	testEqImp(t, imp{cmd: sto}, cd[0])
	testEqImp(t, imp{cmd: lod}, cd[1])
	testEqImp(t, imp{cmd: end}, cd[2])
}

func TestParseFlow(t *testing.T) {
	src := `LSS(mrk) SL
			LST(cll) SSSL     -- 0
			LSL(jmp) SL       -- 1
			LSS(mrk) SSSL
			LTS(jze) SL       -- 2
			LTT(jne) SSSL     -- 3
			LTL(ret)          -- 4
			LLL               -- 5  `
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	testEqInt(t, 6, len(cd)) // not exist, mrk
	testEqImp(t, imp{cmd: cll, arg: 2}, cd[0])
	testEqImp(t, imp{cmd: jmp, arg: 0}, cd[1])
	testEqImp(t, imp{cmd: jze, arg: 0}, cd[2])
	testEqImp(t, imp{cmd: jne, arg: 2}, cd[3])
	testEqImp(t, imp{cmd: ret}, cd[4])
	testEqImp(t, imp{cmd: end}, cd[5])
}
