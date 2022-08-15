package blogrenderer

import (
	"embed"
	"html/template"
	"io"
	"strings"
)

var (
	//go:embed "templates/*"
	postTemplates embed.FS
)

type PostRender struct {
	templ *template.Template
}

func NewPostRender() (*PostRender, error) {
	templ, err := template.ParseFS(postTemplates, "templates/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &PostRender{templ: templ}, nil
}

func (r *PostRender) Render(w io.Writer, p Post) error {
	if err := r.templ.Execute(w, p); err != nil {
		return err
	}

	return nil
}

func (r *PostRender) RenderIndex(w io.Writer, posts []Post) error {
	indexTemplate := `<ol>{{range .}}<li><a href="/post/{{sanitiseTitle .Title}}">{{.Title}}</a></li>{{end}}</ol>`

	templ, err := template.New("index").Funcs(template.FuncMap{
		"sanitiseTitle": func(title string) string {
			return strings.ToLower(strings.Replace(title, " ", "-", -1))
		},
	}).Parse(indexTemplate)
	if err != nil {
		return err
	}

	if err = templ.Execute(w, posts); err != nil {
		return err
	}

	return nil
}
