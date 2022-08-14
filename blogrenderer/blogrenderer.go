package blogrenderer

import (
	"embed"
	"html/template"
	"io"
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
