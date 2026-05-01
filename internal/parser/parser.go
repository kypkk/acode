// Package parser exposes a language-agnostic surface for turning source bytes
// into a syntax tree. The interface is intentionally small; languages are
// added by providing a concrete implementation (see GoParser).
package parser

import sitter "github.com/smacker/go-tree-sitter"

// Tree wraps a tree-sitter parse tree. The underlying *sitter.Tree is embedded
// so callers can use methods like RootNode directly while the wrapper gives us
// a stable seam for adding metadata (file path, language, etc.) later without
// breaking callers.
type Tree struct {
	*sitter.Tree
}

// Parser turns source bytes into a Tree. Each implementation is bound to one
// language.
type Parser interface {
	Parse(src []byte) (*Tree, error)
}
