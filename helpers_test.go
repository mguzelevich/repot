package repot

// https://golang.org/pkg/testing/

import (
	// "bytes"
	// "fmt"
	"testing"
)

func compare(a, b []int) bool {
	if &a == &b {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

type test struct {
	in  string
	out []int
}

func TestParseDeltaString_pos(t *testing.T) {
	testes := []test{
		test{"1", []int{1}},
		test{"1-2", []int{1, 2}},
		test{"1-3", []int{1, 2, 3}},
		test{"1,3-7,12", []int{1, 3, 4, 5, 6, 7, 12}},
	}
	for i, tst := range testes {
		if in, err := ParseRangesString(tst.in); err != nil || !compare(in, tst.out) {
			t.Error(
				"For", "TestParseDeltaString", i,
				"expected", tst.out,
				"got", in,
			)
		}
	}
}

func TestParseDeltaString_neg(t *testing.T) {
	testes := []test{
		//test{"1", []int{1}},
		test{"", []int{1}},
		test{" ", []int{1}},
		test{"-2", []int{1}},
		test{"1-", []int{1}},
		test{",3", []int{1}},
		test{"3,", []int{1}},
	}
	for i, tst := range testes {
		if in, err := ParseRangesString(tst.in); err == nil {
			t.Error(
				"For", "TestParseDeltaString", i,
				"expected", "ERROR",
				"got", in,
			)
		}
	}
}
