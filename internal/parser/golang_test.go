package parser

import "testing"

// Compile-time assertion: GoParser must satisfy Parser.
var _ Parser = (*GoParser)(nil)

func TestGoParser_ParsesValidGoSource(t *testing.T) {
	p := NewGoParser()
	src := []byte("package foo\n\nfunc Bar() {}\n")

	tree, err := p.Parse(src)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if tree == nil {
		t.Fatal("Parse returned nil tree")
	}

	root := tree.RootNode()
	if root == nil {
		t.Fatal("RootNode is nil")
	}
	if got := root.Type(); got != "source_file" {
		t.Fatalf("root type = %q, want %q", got, "source_file")
	}
}
