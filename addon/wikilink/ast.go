package wikilink

import (
	"github.com/yuin/goldmark/ast"
)

var WikilinkKind = ast.NewNodeKind("Wikilink")

type WikilinkNode struct {
	*ast.Image
	IsLink bool // true for link wikilinks, false for image wikilinks
}

func (n *WikilinkNode) Kind() ast.NodeKind {
	return WikilinkKind
}

func (n *WikilinkNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}
