package blogrenderer

import (
	"io"
  "html/template"
)

const (
  postTemplate = `<h1>{{.Title}}</h1>
<p>{{.Description}}</p>
Tags: <ul>{{range .Tags}}<li>{{.}}</li>{{end}}</ul>`
)

func Render(w io.Writer, p Post) error {
  templ, err := template.New("blog").Parse(postTemplate)
	if err != nil {
		return err
	}

  if err := templ.Execute(w, p); err != nil {
    return err
  }

	return nil
}
