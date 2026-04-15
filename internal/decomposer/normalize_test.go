package decomposer

import "testing"

func TestNormalizeFunctionSig(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"@ (a: i32, b: i32) -> i32", "(a: i32, b: i32) -> i32"},
		{"@(a: i32) -> i32", "(a: i32) -> i32"},
		{"@  (a: i32) -> i32", "(a: i32) -> i32"},
		{"@\t(a: i32) -> i32", "(a: i32) -> i32"},
		{"fn (a: i32) -> i32", "(a: i32) -> i32"},
		{"fn(a: i32) -> i32", "(a: i32) -> i32"},
		{"FN (a: i32) -> i32", "(a: i32) -> i32"},
		{"Fn (a: i32) -> i32", "(a: i32) -> i32"},
		{"  @ (a: i32) -> i32", "(a: i32) -> i32"},
		{"\t@ (a: i32) -> i32", "(a: i32) -> i32"},
		{"fn @ (a: i32) -> i32", "(a: i32) -> i32"},
		{"@ fn (a: i32) -> i32", "(a: i32) -> i32"},
		{"(a: i32) -> i32", "(a: i32) -> i32"},
		{"() -> string", "() -> string"},
		{"", ""},
		{"   ", ""},
		{"@", ""},
		{"fn", ""},
		{"FN", ""},
	}
	for _, tc := range cases {
		got := NormalizeFunctionSig(tc.in)
		if got != tc.want {
			t.Errorf("NormalizeFunctionSig(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
