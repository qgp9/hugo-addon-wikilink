package wikilink

import (
	"github.com/gohugoio/hugo/markup/converter"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer"
)

//=======================================================================
// Loader
//=======================================================================

// Load configures wikilink extension and renderer based on the provided config.
// It reads configuration from params.wikilink and appends the extension and renderer options to the provided slices if wikilink is enabled.
func Load(pcfg converter.ProviderConfig, extensions *[]goldmark.Extender, rendererOptions *[]renderer.Option) {
	// Read wikilink config from params.wikilink
	cfg := NewWikilinkConfig(pcfg)
	if !cfg.Enable {
		return
	}

	*extensions = append(*extensions, New(cfg))
	// Note: Renderer registration is now handled in the Extend method
}

// New creates a new Wikilink Goldmark extension.
func New(cfg WikilinkConfig) goldmark.Extender {
	if !cfg.Enable {
		return noOpExtender{}
	}
	return &WikilinkExtension{cfg: cfg}
}
