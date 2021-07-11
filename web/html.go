package web

import (
	"github.com/amery/go-webpack-starter/web/assets"
	"github.com/amery/go-webpack-starter/web/html"
)

func CompileHtml(hashifyAssets bool) (*html.Collection, error) {
	// bind assets to html templates
	h, err := html.Files.Clone()
	if err != nil {
		return nil, err
	}
	h.Funcs(assets.Files.FuncMap(hashifyAssets, "File"))
	// compile templates
	if err := h.Parse(); err != nil {
		return nil, err
	}

	return &h, nil
}
