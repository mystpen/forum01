package handler

import (
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", h.home)

	mux.HandleFunc("/signup", h.signup)

	mux.HandleFunc("/signin", h.signIn)
	mux.HandleFunc("/logout", h.logout)
	mux.HandleFunc("/post/create", h.requireAuth(h.postCreate))
	mux.HandleFunc("/post", h.post)
	mux.HandleFunc("/post/like", h.postLike)
	mux.HandleFunc("/comment/like", h.commentLike)
	mux.HandleFunc("/comment/create", h.commentCreate)
	mux.HandleFunc("/created", h.requireAuth(h.createdPosts))
	mux.HandleFunc("/mylikes", h.requireAuth(h.likedPosts))

	return h.middleware(mux)
}
