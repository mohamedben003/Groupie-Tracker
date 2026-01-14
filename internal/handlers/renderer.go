package handlers

import (
	"html/template"
	"net/http"
)

type TemplateRenderer struct {
	tpls *template.Template
}

func NewTemplateRenderer(glob string) (*TemplateRenderer, error) {
	t, err := template.ParseGlob(glob)
	if err != nil {
		return nil, err
	}
	return &TemplateRenderer{tpls: t}, nil
}

func (r *TemplateRenderer) ExecuteTemplate(w http.ResponseWriter, name string, data any) error {
	return r.tpls.ExecuteTemplate(w, name, data)
}


