package email

import (
	"bytes"
	"embed"
	"fmt"
	htmlTemplate "html/template"
	textTemplate "text/template"
)

//go:embed templates/*.html templates/*.txt
var templatesFS embed.FS

// Renderer holds the parsed HTML and text template trees.
type Renderer struct {
	html *htmlTemplate.Template
	text *textTemplate.Template
}

func NewRenderer() (*Renderer, error) {
	htmlTpl, err := htmlTemplate.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse html templates: %w", err)
	}
	textTpl, err := textTemplate.ParseFS(templatesFS, "templates/*.txt")
	if err != nil {
		return nil, fmt.Errorf("parse text templates: %w", err)
	}
	return &Renderer{html: htmlTpl, text: textTpl}, nil
}

// Render executes both the .html and .txt template for the given base name
// (e.g. "welcome", "password_reset") with data.
func (r *Renderer) Render(name string, data any) (htmlBody, textBody string, err error) {
	var hb, tb bytes.Buffer
	if err := r.html.ExecuteTemplate(&hb, name+".html", data); err != nil {
		return "", "", fmt.Errorf("render html %s: %w", name, err)
	}
	if err := r.text.ExecuteTemplate(&tb, name+".txt", data); err != nil {
		return "", "", fmt.Errorf("render text %s: %w", name, err)
	}
	return hb.String(), tb.String(), nil
}
