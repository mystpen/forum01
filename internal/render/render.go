package render

import (
	"fmt"
	"forum/internal/types"
	"html/template"
	"net/http"
)

type WebPage struct {
	IsLoggedin bool
	Post *types.Post
	Posts []*types.Post
	Errtext string
}

// func NewWebPage(value bool) *WebPage {
// 	return &WebPage{value}
// }

func Render(w http.ResponseWriter, temp string, data WebPage) {

	// fmt.Printf("temp: %v\n", temp)
	templ, err := template.ParseFiles(temp, "ui/html/layout.html")
	if err != nil {
		// ErrorPage(w, r, http.StatusNotFound)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, data)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		// ErrorPage(w, r, http.StatusInternalServerError)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
