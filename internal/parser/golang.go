package parser

import (
	"context"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// GoParser parses Go source.
type GoParser struct {
	p *sitter.Parser
}

// NewGoParser builds a GoParser ready to parse Go source.
func NewGoParser() *GoParser {
	p := sitter.NewParser()
	p.SetLanguage(golang.GetLanguage())
	return &GoParser{p: p}
}

// Parse turns Go source into a Tree.
func (g *GoParser) Parse(src []byte) (*Tree, error) {
	t, err := g.p.ParseCtx(context.Background(), nil, src)
	if err != nil {
		return nil, err
	}
	return &Tree{Tree: t}, nil
}
