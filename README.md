# Wikilink Addon for Hugo

Adds wiki-style links to Hugo.  Use `[[page-name]]` for links and `![[image-name]]` for images. Like Obsidian/Notion syntax.

* Uses https://github.com/qgp9/hugo-addon
* Uses https://pkg.go.dev/go.abhg.dev/goldmark/wikilink parser.

## TL;DR

### Quick setup

```toml
[params.addon.wikilink]
enable = true # Required: explicitly enable the addon
# By default, it uses Hugo's standard link/image render hooks.
```

### Usage

```markdown
[[my-page]]                    # Basic link
[[page|custom text]]           # Custom link text  
![[my-image]]                  # Embed image
![[image|custom alt]]          # Custom alt text
```

## Build

```bash
git clone --recursive --depth=1 --shallow-submodules https://github.com/qgp9/hugo-addon
cd hugo-addon
./bin/patch.sh check 
./bin/patch.sh apply
cp -r addons/wikilink/addon hugo/

cd hugo
go mod tidy
go build -tags wikilink
./hugo version
# hugo v0.148.1-65267cc6ffc4098cd3b8b809ee57b8a0c20f2351+wikilink darwin/arm64 BuildDate=
go test -tags wikilink addon/wikilink/tests/*
```

## Features

* **Wiki-style links**: Use `[[page-name]]` to link to internal pages
* **Custom link text**: Use `[[page-name|custom text]]` for custom link text
* **Image embedding**: Use `![[image-name]]` to embed images
* **Custom alt text**: Use `![[image-name|custom alt]]` for custom alt text
* **Flexible rendering**: Support for both standard link/image hooks and custom wikilink render hooks
* **Build tag protection**: Only compiled when `wikilink` build tag is enabled

## Configuration

Enable the wikilink addon in your Hugo configuration:

```toml
[params.addon.wikilink]
enable = true                # Required: explicitly enable the addon
useStandardLinkHook = true   # Optional: defaults to true
```

### Configuration Options

* `enable` (bool): Enable/disable the wikilink addon (**default: `false`**)
* `useStandardLinkHook` (bool): Use standard link render hooks instead of custom wikilink hooks (**default: `true`**)

**Important**: The addon is **disabled by default** (`enable = false`). You must explicitly enable it in your configuration to use wikilink features.

## Rendering Options

The addon supports two rendering modes:

### Standard Link Hook Mode (`useStandardLinkHook = true`)

Uses Hugo's standard link and image render hooks. This is the **default mode** and works with existing render hook templates.

```sh
layouts/_markup/render-link.html
layouts/_markup/render-image.html
```

### Custom Wikilink Hook Mode (`useStandardLinkHook = false`)

Uses custom render hooks specifically for wikilinks. This mode provides more control over wikilink rendering.

```sh
# Note: must use "_" instead of "-" for the wikilink_link and wikilink_image.
layouts/_markup/render-addon-wikilink_link.html
layouts/_markup/render-addon-wikilink_image.html
```

* A template *"may be"* compatible with standard link/image render hooks.
  * Embedded image render hook:
    [Docs](https://gohugo.io/render-hooks/images/) /
    [Template](https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_markup/render-image.html)
  * Embedded link render hook:
    [Docs](https://gohugo.io/render-hooks/links/) /
    [Template](https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/_markup/render-link.html)

## Build Requirements

The wikilink addon requires the `wikilink` build tag to be enabled:

```bash
go build -tags wikilink
```

Or using mage:

```bash
HUGO_BUILD_TAGS=wikilink mage hugo
```

For testing:

```bash
go test -tags wikilink ./addon/wikilink/tests/*
```

## Dependencies

The addon depends on:

* `go.abhg.dev/goldmark/wikilink` - Core wikilink parsing functionality

## Implementation Details

### Architecture

The wikilink addon is implemented as a Hugo addon that extends Goldmark's markdown processing:

1. **Parser**: Wraps the upstream wikilink parser and converts wikilink nodes to Hugo-compatible AST nodes
2. **Renderer**: Provides custom rendering logic for wikilink nodes
3. **Configuration**: Integrates with Hugo's configuration system
4. **Build Protection**: Only compiled when the `wikilink` build tag is present

### Key Components

* `WikilinkParser`: Parses `[[...]]` and `![[...]]` syntax
* `WikilinkNode`: Custom AST node for wikilinks
* `WikilinkExtension`: Goldmark extension that registers the parser and renderer
* `WikilinkConfig`: Configuration management

### Node Types

* **Link wikilinks**: `[[page-name]]` → `ast.Link` or `WikilinkNode`
* **Image wikilinks**: `![[image-name]]` → `ast.Image` or `WikilinkNode`

## Limitations && Known Issues

* Wikilink render hook may not be 100% compatible with standard link/image render hooks templates.
* `layouts/_markup/render-addon.html` will cause unexpected behavior.  
   Do not use it without good reason.

## A Note from the Developer

1. **No PR intended**: This was built for personal convenience, not as a contribution to Hugo core. The implementation takes some liberties with Hugo's internals that wouldn't pass a proper code review.
1. **Passthrough Extension**: Hugo's Passthrough extension with higher priority may be an easier way to add wikilinks even though it still requires some hacks.
1. **`params.addon.wikilink`**: I didn't want to touch `markup.goldmark.*`.
1. An own wiki-link parser would be better to avoid multiple conversions between AST nodes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
