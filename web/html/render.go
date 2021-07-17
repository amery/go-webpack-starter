package html

import (
	"log"
	"net/http"

	"go.sancus.dev/core/context"
	"go.sancus.dev/core/errors"
	"go.sancus.dev/file2go/html"
	"go.sancus.dev/web"
)

type Collection = html.Collection

// Middleware that attaches a HTML Template context to the request
func Middleware(h *html.Collection) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := WithHtmlContext(r.Context(), h)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// Returns a Renderer HTML Template attached to the given data
func View(r *http.Request, name string, data interface{}) web.Renderer {

	if t := HtmlContext(r.Context()); t == nil {
		// Not Found
		log.Fatal(ErrContextNotFound)
	} else if v, err := t.View(name, data); err == nil {
		// View Ready
		return v
	} else {
		// Template Not found
		log.Fatal(err)
	}
	return nil
}

// HtmlContext returns the HTML Template context attached to a context
func HtmlContext(ctx context.Context) *html.Collection {
	if t, ok := ctx.Value(collectionCtxKey).(*html.Collection); ok {
		return t
	} else {
		return nil
	}
}

// WithHtmlContext returns a new http.Request Context with a given
// HTML Template context attached if given
func WithHtmlContext(ctx context.Context, t *html.Collection) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if t != nil {
		ctx = context.WithValue(ctx, collectionCtxKey, t)
	}
	return ctx
}

var (
	// collectionCtxKey references a collection of html templates
	collectionCtxKey = context.NewContextKey("HtmlTemplatesContext")

	// Error raised by View() when the HtmlContext isn't available
	ErrContextNotFound = errors.New("%s not found in %s", "HtmlContext", "http.Request")
)
