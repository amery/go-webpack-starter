package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/amery/go-webpack-starter/web/assets"
	"github.com/amery/go-webpack-starter/web/html"
)

type Router struct {
	http.Handler

	HashifyAssets bool
}

func (c *Router) Compile() error {
	// bind assets to html templates
	h, err := html.Files.Clone()
	if err != nil {
		return err
	}
	h.Funcs(assets.Files.FuncMap(c.HashifyAssets, "File"))

	// compile templates
	if err := h.Parse(); err != nil {
		return err
	}

	// and compose then router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(assets.Middleware(c.HashifyAssets))
	r.Use(html.Middleware(&h))
	r.MethodFunc("GET", "/", HandleIndex)

	c.Handler = r
	return nil
}
