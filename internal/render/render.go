package render

import (
	"bytes"
	"html/template"
	"net/http"
)

var templates *template.Template

func Init(glob string) error {
	t, err := template.ParseGlob(glob)
	if err != nil {
		return err
	}
	templates = t
	return nil
}

func HTML(w http.ResponseWriter, name string, data any) error {
	
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		return err
	}
	_, err := w.Write(buf.Bytes())
	return err
}


