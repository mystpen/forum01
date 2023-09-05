package handler

import (
	"fmt"
	"net/http"

	"forum/internal/cookies"
	"forum/internal/render"
	"forum/internal/types"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
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
	switch r.Method {
	case http.MethodGet:

		var categories []string
		posts := []*types.Post{}
		var err error

		query := r.URL.Query()
		categories = append(categories, query["category"]...)

		// for _, value := range query["category"] {
		// 	categories = append(categories, value)
		// }
		if len(categories) == 0 {
			posts, err = h.service.PostService.GetAllPosts()
			if err != nil {
				fmt.Printf("err: %v\n", err)
				ErrorPage(w, r, http.StatusInternalServerError)
				return
			}
		} else {
			posts, err = h.service.PostService.Filter(categories)
			// GetPostsByCategory(categories)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				ErrorPage(w, r, http.StatusInternalServerError)
				return
			}
		}

		render.Render(w, "ui/html/index.html", render.WebPage{
			IsLoggedin: data,
			Posts:      posts,
		})

	default:
		ErrorPage(w, r, http.StatusMethodNotAllowed)
	}
}
