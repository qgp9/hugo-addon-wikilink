package wikilink

import (
	"github.com/gohugoio/hugo/addon/common"
	"github.com/gohugoio/hugo/markup/converter"
)

type WikilinkConfig struct {
	Enable              bool
	UseStandardLinkHook bool // If true, use standard link render hook; if false, use custom wikilink render hook
}

func NewWikilinkConfig(pcfg converter.ProviderConfig) WikilinkConfig {
	wcfg := WikilinkConfig{
		Enable:              false,
		UseStandardLinkHook: true, // Default to true (use standard link render hook)
	}
	common.GetAddonConfig(pcfg.Conf.GetConfig(), "wikilink", &wcfg)
	return wcfg
}
