package analyzer

import "testing"

func TestPackageName(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want string
	}{
		{"simple", "package x\n", "x"},
		{"with body", "package mypkg\n\nfunc F() {}\n", "mypkg"},
		{"with leading comment", "// header\npackage analyzer\n", "analyzer"},
		{"empty source", "", ""},
		{"no package clause", "// only a comment\n", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tree, src := parseGo(t, c.src)
			if got := PackageName(tree, src); got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}
