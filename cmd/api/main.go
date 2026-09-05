package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Handle(
		"/static/*",
		http.StripPrefix(
			"/static/",
			http.FileServer(
				http.Dir("./web/static"),
			),
		),
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(
			w,
			r,
			"./web/static/index.html",
		)
	})

	fmt.Println("Bukuin is running on http://localhost:8080")

	http.ListenAndServe(":8080", r)
}
