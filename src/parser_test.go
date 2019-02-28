package main

import (
	"testing"
)

const errorFormat = "expected:%v, actual:%v"

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

func TestParsePush(t *testing.T) {
	src := "SS(push),STL(+1),SS(push),TTSTL(-101),LLL"
	cd, err := parse(mkSource(src))
	if err != nil {
		t.Error(err)
	}
	if 3 != len(cd) {
		t.Errorf(errorFormat, 3, len(cd))
	}
	if (imp{cmd: psh, arg: 1}) != cd[0] {
		t.Errorf(errorFormat, imp{cmd: psh, arg: 1}, cd[0])
	}
	if (imp{cmd: psh, arg: -5}) != cd[1] {
		t.Errorf(errorFormat, imp{cmd: psh, arg: -5}, cd[1])
	}
	if (imp{cmd: end}) != cd[2] {
		t.Errorf(errorFormat, imp{cmd: end}, cd[2])
	}
}
