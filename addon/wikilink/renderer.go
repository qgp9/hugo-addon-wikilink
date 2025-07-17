package wikilink

import (
	"github.com/gohugoio/hugo/addon/common"
	commonrenderer "github.com/gohugoio/hugo/addon/common/renderer"
	"github.com/gohugoio/hugo/common/types/hstring"
	"github.com/gohugoio/hugo/markup/addon/attributes"
	"github.com/gohugoio/hugo/markup/goldmark/addon/render"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type htmlRenderer struct {
}

func NewHTMLRenderer() renderer.NodeRenderer {
	return &htmlRenderer{}
}

func (r *htmlRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(WikilinkKind, r.renderWikilink) // WikilinkNode kind
}

// renderWikilink handles WikilinkNode (converted from wikilinks)
func (r *htmlRenderer) renderWikilink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	wikilinkNode := node.(*WikilinkNode)

	// renderID should be word-safe(`[a-zA-Z0-9_]+`: no "-")
	// DO NOT "wikilink-link" or "wikilink-image"
	var renderId string
	if wikilinkNode.IsLink {
		renderId = "wikilink_link"
	} else {
		renderId = "wikilink_image"
	}

	wlRenderer := commonrenderer.GetAddonRenderer(w, renderId)
	if wlRenderer == nil {
		// No wikilink renderer found, skip rendering
		return ast.WalkContinue, nil
	}

	ctx, ok := w.(*render.Context)
	if !ok {
		return ast.WalkContinue, nil
	}

	if entering {
		// Store the current pos so we can capture the rendered text.
		ctx.PushPos(ctx.Buffer.Len())
		return ast.WalkContinue, nil
	}

	text := ctx.PopRenderedString()

	// Use different ordinal based on whether it's a link or image
	var ordinal int
	if wikilinkNode.IsLink {
		ordinal = ctx.GetAndIncrementOrdinal(ast.KindLink)
	} else {
		ordinal = ctx.GetAndIncrementOrdinal(ast.KindImage)
	}

	page, pageInner := render.GetPageAndPageInner(ctx)

	// Create wikilink context using existing structures from render_hooks.go
	wctx := &wikilinkContext{
		ImageLinkContext: common.NewImageLinkContext(
			common.NewLinkContext(
				page,
				pageInner,
				string(wikilinkNode.Image.Destination),
				"", // wikilink does not have a title
				hstring.HTML(text),
				render.TextPlain(wikilinkNode.Image, source),
				attributes.Empty,
			),
			ordinal,
			false, // WikilinkNode does not have an Embed field
		),
		isLink: wikilinkNode.IsLink,
	}

	if err := wlRenderer.RenderAddon(ctx.RenderContext().Ctx, w, wctx); err != nil {
		return ast.WalkStop, err
	}

	return ast.WalkContinue, nil
}

// wikilinkContext combines linkContext and imageLinkContext to satisfy WikilinkContext interface
type wikilinkContext struct {
	common.ImageLinkContext
	isLink bool
}

func (w *wikilinkContext) IsLink() bool {
	return w.isLink
}
