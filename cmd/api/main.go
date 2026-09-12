package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/admin"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/auth"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/book"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/cart"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/database"
	"github.com/ranggawaridat/bukuin-digital-bookstore/internal/order"
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

	cartRepository := cart.NewRepository(db)
	cartHandler := cart.NewHandler(
		cartRepository,
	)

	orderRepository := order.NewRepository(db)
	orderHandler := order.NewHandler(
		orderRepository,
	)

	adminHandler := admin.NewHandler(db)

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
		"/books/{id}",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/book.html",
			)
		},
	)

	r.Get(
		"/cart",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/cart.html",
			)
		},
	)

	r.Get(
		"/orders",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/orders.html",
			)
		},
	)

	r.Get(
		"/orders/{id}",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/order.html",
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
		"/profile",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/profile.html",
			)
		},
	)

	r.Get(
		"/admin",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			http.ServeFile(
				w,
				r,
				"./web/static/admin.html",
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

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAuth(next.ServeHTTP)(w, r)
		})
	}).Get(
		"/api/cart",
		cartHandler.GetCart,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAuth(next.ServeHTTP)(w, r)
		})
	}).Post(
		"/api/cart/items",
		cartHandler.AddItem,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAuth(next.ServeHTTP)(w, r)
		})
	}).Put(
		"/api/cart/items/{id}",
		cartHandler.UpdateItem,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAuth(next.ServeHTTP)(w, r)
		})
	}).Delete(
		"/api/cart/items/{id}",
		cartHandler.DeleteItem,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAuth(next.ServeHTTP)(w, r)
		})
	}).Post(
		"/api/orders/checkout",
		orderHandler.Checkout,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAuth(next.ServeHTTP)(w, r)
		})
	}).Get(
		"/api/orders",
		orderHandler.GetOrders,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAuth(next.ServeHTTP)(w, r)
		})
	}).Get(
		"/api/orders/{id}",
		orderHandler.GetOrderByID,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAdmin(next.ServeHTTP)(w, r)
		})
	}).Get(
		"/api/admin/stats",
		adminHandler.Stats,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAdmin(next.ServeHTTP)(w, r)
		})
	}).Get(
		"/api/admin/orders",
		adminHandler.Orders,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAdmin(next.ServeHTTP)(w, r)
		})
	}).Post(
		"/api/admin/books",
		bookHandler.CreateBook,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAdmin(next.ServeHTTP)(w, r)
		})
	}).Put(
		"/api/admin/books/{id}",
		bookHandler.UpdateBook,
	)

	r.With(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHandler.RequireAdmin(next.ServeHTTP)(w, r)
		})
	}).Delete(
		"/api/admin/books/{id}",
		bookHandler.DeleteBook,
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
