//go:build wikilink

package wikilink_test

import (
	"testing"

	"github.com/gohugoio/hugo/hugolib"
)

// Test cases for simple render hooks
func TestWikilink_Simple_Link_UseStandardLinkHook_True(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"link", renderWikilinkLinkSimple,
		"This is a [[my-page]].\nThis is a [[another-page|custom text]].",
		"<p>This is a LINK-SIMPLE: my-page | my-page.\nThis is a LINK-SIMPLE: another-page | custom text.</p>",
	)
}

func TestWikilink_Simple_Link_UseStandardLinkHook_False(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_link", renderWikilinkLinkSimple,
		"This is a [[my-page]].\nThis is a [[another-page|custom text]].",
		"<p>This is a LINK-SIMPLE: my-page | my-page.\nThis is a LINK-SIMPLE: another-page | custom text.</p>",
	)
}

func TestWikilink_Simple_Image_UseStandardLinkHook_True(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"image", renderWikilinkImageSimple,
		"This is an image ![[my-image]].\nThis is another ![[another-image|custom alt]].",
		"<p>This is an image IMAGE-SIMPLE: my-image | .\nThis is another IMAGE-SIMPLE: another-image | custom alt.</p>",
	)
}

func TestWikilink_Simple_Image_UseStandardLinkHook_False(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_image", renderWikilinkImageSimple,
		"This is an image ![[my-image]].\nThis is another ![[another-image|custom alt]].",
		"<p>This is an image IMAGE-SIMPLE: my-image | .\nThis is another IMAGE-SIMPLE: another-image | custom alt.</p>",
	)
}
func TestWikilink_Standard_Link_UseStandardLinkHook_True(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"link", renderWikilinkLinkStandard,
		"[[my-page]].[[another-page|custom text]].",
		`<p><a href="/my-page/" class="render-wikilink-link">my-page</a>.<a href="/another-page/" class="render-wikilink-link">custom text</a>.</p>`,
	)
}

func TestWikilink_Standard_Link_UseStandardLinkHook_False(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_link", renderWikilinkLinkStandard,
		"[[my-page]].[[another-page|custom text]].",
		`<p><a href="/my-page/" class="render-wikilink-link">my-page</a>.<a href="/another-page/" class="render-wikilink-link">custom text</a>.</p>`,
	)
}

func TestWikilink_Standard_Image_UseStandardLinkHook_True(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"image", renderWikilinkImageStandard,
		"This is an image ![[my-image]].This is another ![[another-image|custom alt]].",
		`<p>This is an image <img src="my-image" alt="" class="render-image2">.This is another <img src="another-image" alt="custom alt" class="render-image2">.</p>`,
	)
}

func TestWikilink_Standard_Image_UseStandardLinkHook_False(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_image", renderWikilinkImageStandard,
		"This is an image ![[my-image]].This is another ![[another-image|custom alt]].",
		`<p>This is an image <img src="my-image" alt="" class="render-image2">.This is another <img src="another-image" alt="custom alt" class="render-image2">.</p>`,
	)
}

// Test configurations
const (
	config_UseStandardLinkHook_true = `
[params.addon.wikilink]
enable = true
useStandardLinkHook = true
`

	config_UseStandardLinkHook_false = `
[params.addon.wikilink]
enable = true
useStandardLinkHook = false
`
)

// Common content files
const (
	contentBasic = `
-- layouts/_default/single.html --
{{ .Content }}

-- layouts/_default/list.html --
{{ .Content }}

-- content/my-page.md --
---
title: "My Page"
---
This is my page content.

-- content/another-page.md --
---
title: "Another Page"
---
This is another page content.

`

	// Render hook templates
	renderWikilinkLinkSimple = `
LINK-SIMPLE: {{ .Destination }} | {{ .Text }}
`
	renderWikilinkImageSimple = `
IMAGE-SIMPLE: {{ .Destination }} | {{ .Text }}
`

	renderWikilinkLinkStandard = `
{{- $url := .Destination -}}
{{- $text := .Text -}}

{{- if hasPrefix $url "http" -}}
  <a href="{{- $url -}}" class="render-link2" target="_blank" rel="noopener">{{- $text -}}</a>
{{- else -}}
  {{- $page := $.PageInner.GetPage $url -}}
  {{- if $page -}}
    <a href="{{- $page.RelPermalink -}}" class="render-wikilink-link">{{- $text -}}</a>
  {{- else -}}
    <a href="{{- $url -}}" class="render-wikilink-link">{{- $text -}}</a>
  {{- end -}}
{{- end -}}
`

	renderWikilinkImageStandard = `
{{- $u := urls.Parse .Destination -}}
{{- $src := $u.String -}}
{{- if not $u.IsAbs -}}
  {{- $path := strings.TrimPrefix "./" $u.Path -}}
  {{- with or (.PageInner.Resources.Get $path) (resources.Get $path) -}}
    {{- $src = .RelPermalink -}}
    {{- with $u.RawQuery -}}
      {{- $src = printf "%s?%s" $src . -}}
    {{- end -}}
    {{- with $u.Fragment -}}
      {{- $src = printf "%s#%s" $src . -}}
    {{- end -}}
  {{- end -}}
{{- end -}}
<img src="{{ $src }}" alt="{{ .PlainText }}" class="render-image2"
  {{- with .Title }} title="{{ . }}" {{- end -}}
  {{- range $k, $v := .Attributes -}}
    {{- if $v -}}
      {{- printf " %s=%q" $k ($v | transform.HTMLEscape) | safeHTMLAttr -}}
    {{- end -}}
  {{- end -}}
>
{{- /**/ -}}
`
)

// Helper function to build test files
func buildTestFiles(config, p1Content, hookType, renderHook string, additionalFiles ...string) string {
	files := `-- config.toml --` + config + `
-- content/docs/p1/_index.md --
---
title: "p1"
---
` + p1Content + contentBasic

	if hookType != "" {
		files += "\n-- layouts/_default/_markup/render-" + hookType + ".html --" + renderHook
	}

	// Add additional files
	for i := 0; i < len(additionalFiles); i += 2 {
		if i+1 < len(additionalFiles) {
			files += "\n-- " + additionalFiles[i] + " --\n" + additionalFiles[i+1]
		}
	}

	return files
}

// Test builder function
func build_test(t *testing.T, config, hookType, renderHook, input, expected string, additionalFiles ...string) *hugolib.IntegrationTestBuilder {
	t.Parallel()

	files := buildTestFiles(config, input, hookType, renderHook, additionalFiles...)

	b := hugolib.Test(t, files)

	// Debug: Print actual HTML output
	// fmt.Println("=== Generated HTML ===")
	// fmt.Println(b.FileContent("public/docs/p1/index.html"))
	// fmt.Println("=== End HTML ===")

	b.AssertFileContent("public/docs/p1/index.html", expected)

	// If a custom render hook is provided, assert that the template file exists
	if renderHook != "" {
		b.AssertFileExists("layouts/_default/_markup/render-"+hookType+".html", true)
	}

	return b
}

// 존재하지 않는 페이지에 대한 위키링크
func TestWikilink_NonExistentPage(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_link", renderWikilinkLinkStandard,
		"This is a [[no-such-page]].",
		`<p>This is a <a href="no-such-page" class="render-wikilink-link">no-such-page</a>.</p>`,
	)
}

// 존재하는 페이지에 대한 위키링크 (useStandardLinkHook = true)
func TestWikilink_ExistingPageUseStandardTrue(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"link", renderWikilinkLinkStandard,
		"This is a [[my-page]].\nThis is a [[another-page|custom text]].",
		`<p>This is a <a href="/my-page/" class="render-wikilink-link">my-page</a>.
This is a <a href="/another-page/" class="render-wikilink-link">custom text</a>.</p>`,
	)
}

// 존재하는 페이지에 대한 위키링크 (useStandardLinkHook = false)
func TestWikilink_ExistingPageUseStandardFalse(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_link", renderWikilinkLinkStandard,
		"This is a [[my-page]].\nThis is a [[another-page|custom text]].",
		`<p>This is a <a href="/my-page/" class="render-wikilink-link">my-page</a>.
This is a <a href="/another-page/" class="render-wikilink-link">custom text</a>.</p>`,
	)
}

// 존재하지 않는 페이지에 대한 위키링크 (useStandardLinkHook = true)
func TestWikilink_NonExistentPageUseStandardTrue(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"link", renderWikilinkLinkStandard,
		"This is a [[no-such-page]].",
		`<p>This is a <a href="no-such-page" class="render-wikilink-link">no-such-page</a>.</p>`,
	)
}

// 존재하지 않는 페이지에 대한 위키링크 (useStandardLinkHook = false)
func TestWikilink_NonExistentPageUseStandardFalse(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_link", renderWikilinkLinkStandard,
		"This is a [[no-such-page]].",
		`<p>This is a <a href="no-such-page" class="render-wikilink-link">no-such-page</a>.</p>`,
	)
}

// Page Bundle 이미지 - 존재하는 이미지 (useStandardLinkHook = true)
func TestWikilink_ExistingImageBundleUseStandardTrue(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"image", renderWikilinkImageStandard,
		"This is an image ![[my-image]].\nThis is another ![[another-image|custom alt]].",
		`<p>This is an image <img src="my-image" alt="" class="render-image2">.
This is another <img src="another-image" alt="custom alt" class="render-image2">.</p>`,
		"content/docs/p1/my-image.png", "\niVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		"content/docs/p1/another-image.jpg", "\n/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwA/8A==",
	)
}

// Page Bundle 이미지 - 존재하는 이미지 (useStandardLinkHook = false)
func TestWikilink_ExistingImageBundleUseStandardFalse(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_image", renderWikilinkImageStandard,
		"This is an image ![[my-image]].\nThis is another ![[another-image|custom alt]].",
		`<p>This is an image <img src="my-image" alt="" class="render-image2">.
This is another <img src="another-image" alt="custom alt" class="render-image2">.</p>`,
		"content/docs/p1/my-image.png", "\niVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		"content/docs/p1/another-image.jpg", "\n/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwA/8A==",
	)
}

// Page Bundle 이미지 - 존재하지 않는 이미지 (useStandardLinkHook = true)
func TestWikilink_NonExistentImageBundleUseStandardTrue(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"image", renderWikilinkImageStandard,
		"This is an image ![[no-such-image]].\nThis is another ![[another-missing|custom alt]].",
		`<p>This is an image <img src="no-such-image" alt="" class="render-image2">.
This is another <img src="another-missing" alt="custom alt" class="render-image2">.</p>`,
	)
}

// Page Bundle 이미지 - 존재하지 않는 이미지 (useStandardLinkHook = false)
func TestWikilink_NonExistentImageBundleUseStandardFalse(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_image", renderWikilinkImageStandard,
		"This is an image ![[no-such-image]].\nThis is another ![[another-missing|custom alt]].",
		`<p>This is an image <img src="no-such-image" alt="" class="render-image2">.
This is another <img src="another-missing" alt="custom alt" class="render-image2">.</p>`,
	)
}

// Assets 이미지 - 존재하는 이미지 (useStandardLinkHook = true)
func TestWikilink_ExistingImageAssetsUseStandardTrue(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"image", renderWikilinkImageStandard,
		"This is an image ![[/images/my-image]].\nThis is another ![[/images/another-image|custom alt]].",
		`<p>This is an image <img src="/images/my-image" alt="" class="render-image2">.
This is another <img src="/images/another-image" alt="custom alt" class="render-image2">.</p>`,
		"assets/images/my-image.png", "\niVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		"assets/images/another-image.jpg", "\n/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwA/8A==",
	)
}

// Assets 이미지 - 존재하는 이미지 (useStandardLinkHook = false)
func TestWikilink_ExistingImageAssetsUseStandardFalse(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_image", renderWikilinkImageStandard,
		"This is an image ![[/images/my-image]].\nThis is another ![[/images/another-image|custom alt]].",
		`<p>This is an image <img src="/images/my-image" alt="" class="render-image2">.
This is another <img src="/images/another-image" alt="custom alt" class="render-image2">.</p>`,
		"assets/images/my-image.png", "\niVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==",
		"assets/images/another-image.jpg", "\n/9j/4AAQSkZJRgABAQEAYABgAAD/2wBDAAEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/2wBDAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQH/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAAAAX/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwA/8A==",
	)
}

// Assets 이미지 - 존재하지 않는 이미지 (useStandardLinkHook = true)
func TestWikilink_NonExistentImageAssetsUseStandardTrue(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_true,
		"image", renderWikilinkImageStandard,
		"This is an image ![[/images/no-such-image]].\nThis is another ![[/images/another-missing|custom alt]].",
		`<p>This is an image <img src="/images/no-such-image" alt="" class="render-image2">.
This is another <img src="/images/another-missing" alt="custom alt" class="render-image2">.</p>`,
	)
}

// Assets 이미지 - 존재하지 않는 이미지 (useStandardLinkHook = false)
func TestWikilink_NonExistentImageAssetsUseStandardFalse(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_image", renderWikilinkImageStandard,
		"This is an image ![[/images/no-such-image]].\nThis is another ![[/images/another-missing|custom alt]].",
		`<p>This is an image <img src="/images/no-such-image" alt="" class="render-image2">.
This is another <img src="/images/another-missing" alt="custom alt" class="render-image2">.</p>`,
	)
}

// 이미지 alt가 비어있는 경우
func TestWikilink_ImageEmptyAlt(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_image", renderWikilinkImageSimple,
		"This is an image ![[img|]].",
		"<p>This is an image ![[img|]].</p>",
	)
}

// UTF-8 문자 (한글, 이모지)
func TestWikilink_UTF8Characters(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_link", renderWikilinkLinkSimple,
		"이것은 [[한국어-페이지]]입니다. 이것은 [[emoji-😀|😀 이모지]]입니다.",
		"<p>이것은 LINK-SIMPLE: 한국어-페이지 | 한국어-페이지입니다. 이것은 LINK-SIMPLE: emoji-😀 | 😀 이모지입니다.</p>",
	)
}

// 텍스트 파일을 추가하고 링크 렌더링을 테스트
func TestWikilink_LinkToTxtFile(t *testing.T) {
	build_test(t, config_UseStandardLinkHook_false,
		"addon-wikilink_link", renderWikilinkLinkStandard,
		"This is a [[myfile.txt]].",
		`<p>This is a <a href="myfile.txt" class="render-wikilink-link">myfile.txt</a>.</p>`,
		"content/docs/p1/myfile.txt", "\nhello txt",
	)
}
