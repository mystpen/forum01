package handler

import (
	"net/http"
	"strconv"

	"forum/internal/cookies"
	"forum/internal/render"
)

func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(w, r, http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || (id < 1) {
		
		ErrorPage(w, r, http.StatusNotFound)
		return
	}
	post, err := h.service.PostService.GetPostByID(id)
	if err != nil {
		ErrorPage(w, r, http.StatusNotFound)
		return
	}
	_, errC := cookies.GetCookie(r)

	var data bool
	if errC != nil {
		data = false
	} else {
		data = true
	}

	render.Render(w, "ui/html/post.html", render.WebPage{
		IsLoggedin: data,
		Post:       post,
	})
}
