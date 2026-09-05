package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/database"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

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

	log.Fatal(
		http.ListenAndServe(":8080", r),
	)
}
