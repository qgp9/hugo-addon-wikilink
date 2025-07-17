package wikilink

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type WikilinkExtension struct {
	cfg WikilinkConfig
}

func (we *WikilinkExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewWikilinkParser(we.cfg.UseStandardLinkHook), 199),
		),
	)

	if !we.cfg.UseStandardLinkHook {
		// Register the custom renderer for wikilinks
		m.Renderer().AddOptions(renderer.WithNodeRenderers(
			util.Prioritized(NewHTMLRenderer(), 198),
		))
	}
}

type noOpExtender struct{}

func (noOpExtender) Extend(m goldmark.Markdown) {}
