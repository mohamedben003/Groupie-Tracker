package render

import (
	"fmt"
	"net/http"

	types "grouping_tracker/internal/types"
)

func RenderError(w http.ResponseWriter, code int, title, message string) {
	w.WriteHeader(code)
	data := types.ErrorPageData{
		Code:    code,
		Title:   title,
		Message: message,
	}

	err := types.Templates.ExecuteTemplate(w, "error.html", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%d %s", code, title), code)
	}
}

func Render404(w http.ResponseWriter) {
	RenderError(w, http.StatusNotFound, "Page Not Found", "Got lost? It seems that the page you're looking for doesn't exist.")
}

func Render500(w http.ResponseWriter) {
	RenderError(w, http.StatusInternalServerError, "Internal Server Error", "Something went wrong on our end. Please try again later.")
}
