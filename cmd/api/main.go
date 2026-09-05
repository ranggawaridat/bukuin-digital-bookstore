package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/auth"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/book"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/database"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(
			"failed to connect database:",
			err,
		)
	}

	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(
			"failed to run migration:",
			err,
		)
	}

	bookRepository := book.NewRepository(db)
	bookHandler := book.NewHandler(
		bookRepository,
	)

	authRepository := auth.NewRepository(db)
	authHandler := auth.NewHandler(
		authRepository,
	)

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

	r.Get(
		"/",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/index.html",
			)
		},
	)

	r.Get(
		"/books",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/books.html",
			)
		},
	)

	r.Get(
		"/login",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/login.html",
			)
		},
	)

	r.Get(
		"/register",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/register.html",
			)
		},
	)

	r.Get(
		"/api/books",
		bookHandler.GetBooks,
	)

	r.Get(
		"/api/books/{id}",
		bookHandler.GetBookByID,
	)

	r.Post(
		"/api/auth/register",
		authHandler.Register,
	)

	r.Post(
		"/api/auth/login",
		authHandler.Login,
	)

	r.Post(
		"/api/auth/logout",
		authHandler.Logout,
	)

	r.Get(
		"/api/auth/me",
		authHandler.Me,
	)

	fmt.Println(
		"Bukuin is running on http://localhost:8080",
	)

	log.Fatal(
		http.ListenAndServe(
			":8080",
			r,
		),
	)
}
