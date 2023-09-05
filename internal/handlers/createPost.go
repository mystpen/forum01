package handler

import (
	"fmt"
	"forum/internal/render"
	"forum/internal/types"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func (h *Handler) postCreate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		ErrorPage(w, r, http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
	} else if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(20 << 20) // max size 20MB
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = r.ParseForm()
		if err != nil {
			fmt.Println(err)
			return
		}

		author := h.getUserFromContext(r)

		categoriesForm := r.PostForm["categories"]

		post := &types.CreatePost{
			AuthorId:   author.Id,
			AuthorName: author.Username,
			Title:      strings.TrimSpace(r.Form.Get("title")),
			Content:    strings.TrimSpace(r.Form.Get("text")),
			Categories: categoriesForm,
		}

		if len(post.Title) == 0 || len(post.Content) == 0 {
			var ErrText string = "Please fill out all fields"
			w.WriteHeader(http.StatusBadRequest)
			render.Render(w, "ui/html/createpost.html", render.WebPage{
				Errtext: ErrText,
			})
			return
		}

		// image upload ///////////////////////////////////////////////

		file, header, err := r.FormFile("image")
		if err != nil {
			if err != http.ErrMissingFile {
				fmt.Println(err)
				ErrorPage(w, r, http.StatusInternalServerError)
				return

			}
			post.ImageData = nil
		} else {
			ext := filepath.Ext(header.Filename)
			allowedExts := map[string]bool{".png": true, ".gif": true, ".jpeg": true, ".jpg": true}
			if !allowedExts[ext] {
				var ErrText string = "Invalid file format"
				w.WriteHeader(http.StatusBadRequest)
				render.Render(w, "ui/html/createpost.html", render.WebPage{
					Errtext: ErrText,
				})
				return
			}
			defer file.Close()

			if header.Size > 20 * 1024 * 1024{
				var ErrText string = "File must be less than 20MB"
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				render.Render(w, "ui/html/createpost.html", render.WebPage{
					Errtext: ErrText,
				})
				return
			}
			post.ImageFormat = ext[1:]
			post.ImageData, err = ioutil.ReadAll(file)
			if err != nil {
				ErrorPage(w, r, http.StatusInternalServerError)
				return
			}

		}

		_, err = h.service.PostService.CreateNewPost(post)
		if err != nil {
			fmt.Println(err)
			ErrorPage(w, r, http.StatusInternalServerError)
			return
		}

		// http.Redirect(w, r, fmt.Sprintf("/post?id=%v", id), http.StatusSeeOther)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		ErrorPage(w, r, http.StatusMethodNotAllowed)
		return
	}
	render.Render(w, "ui/html/createpost.html", render.WebPage{})
}
