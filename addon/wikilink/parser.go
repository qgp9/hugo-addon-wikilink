package wikilink

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/wikilink"
)

// TODO: copy or own parser? Too many conversions from wikilink.Node

// WikilinkParser wraps the upstream wikilink.Parser and converts wikilink.Node to ast.Link/ast.Image or WikilinkNode
type WikilinkParser struct {
	upstreamParser      parser.InlineParser
	useStandardLinkHook bool
}

// NewWikilinkParser creates a new wikilink parser that wraps the upstream parser
func NewWikilinkParser(useStandardLinkHook bool) parser.InlineParser {
	return &WikilinkParser{
		upstreamParser:      &wikilink.Parser{},
		useStandardLinkHook: useStandardLinkHook,
	}
}

// Trigger returns characters that trigger this parser
func (p *WikilinkParser) Trigger() []byte {
	return p.upstreamParser.Trigger()
}

// Parse parses wikilinks using upstream parser and converts wikilink.Node to ast.Link/ast.Image or WikilinkNode
func (p *WikilinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	// Use upstream parser to get wikilink.Node
	node := p.upstreamParser.Parse(parent, block, pc)

	wikilinkNode, ok := node.(*wikilink.Node)
	if !ok {
		return node // Return as-is if not a wikilink node
	}

	destination := string(wikilinkNode.Target)
	text := destination

	// First child is always a text node and not empty.
	// Filled with destination unless custom text is provided.
	textNode, ok := wikilinkNode.FirstChild().(*ast.Text)
	if ok {
		text = string(textNode.Segment.Value(block.Source()))
	}
	hasText := text != destination
	isLink := !wikilinkNode.Embed

	// Create ast.Link node first
	link := ast.NewLink()
	link.Destination = wikilinkNode.Target
	link.Title = nil

	var ret_node ast.Node = link

	// Link has text always.
	// Image has text only if custom text is provided.
	if isLink || hasText {
		link.AppendChild(link, textNode)
	}

	// only standard link render hook use link node.
	// wikilink node is based on an image even if it's a link.
	if !isLink || !p.useStandardLinkHook {
		ret_node = ast.NewImage(link)
	}

	// Wikilink render hook use WikilinkNode instead of ast.Link or ast.Image.
	if !p.useStandardLinkHook {
		ret_node = &WikilinkNode{
			Image:  ret_node.(*ast.Image),
			IsLink: isLink,
		}
	}

	return ret_node
}
