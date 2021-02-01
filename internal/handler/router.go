package handler

import (
	"gfile/pkg/middleware"
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (c *Configuration) Route() http.Handler {
	// middleware
	m0 := alice.New(middleware.RecoverPanic, middleware.LogRequest, middleware.DefaultHeaders)
	m1 := alice.New(c.Session.Enable, c.authenticate, c.authorize)

	r := pat.New()

	// static files
	r.Get("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	r.Get("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/template/js"))))

	r.Get("/", m1.Append(c.requireAuthenticatedUser).ThenFunc(c.File.Index(c)))

	r.Get("/users/login", m1.ThenFunc(c.User.LoginView(c)))
	r.Post("/users/login", m1.Then(http.HandlerFunc(c.User.Login(c))))
	r.Post("/users/logout", m1.Then(http.HandlerFunc(c.User.Logout(c))))

	r.Get("/users", m1.ThenFunc(c.User.Index(c)))
	r.Get("/users/new", m1.ThenFunc(c.User.AddView(c)))
	r.Post("/users", m1.ThenFunc(c.User.Add(c)))
	r.Del("/users/:id", m1.Then(http.HandlerFunc(c.User.Delete)))
	r.Get("/users/:id", m1.ThenFunc(c.User.EditView(c)))
	r.Put("/users/:id", m1.Then(http.HandlerFunc(c.User.Edit)))

	r.Get("/roles", m1.ThenFunc(c.Role.Index(c)))
	r.Get("/roles/new", m1.ThenFunc(c.Role.AddView(c)))
	r.Post("/roles", m1.ThenFunc(c.Role.Add(c)))
	r.Del("/roles/:id", m1.Then(http.HandlerFunc(c.Role.Delete)))
	r.Get("/roles/:id", m1.ThenFunc(c.Role.EditView(c)))
	r.Put("/roles/:id", m1.Then(http.HandlerFunc(c.Role.Edit)))

	r.Get("/auth", m1.ThenFunc(c.Casbin.Index(c)))
	r.Get("/auth/new", m1.ThenFunc(c.Casbin.AddView(c)))
	r.Post("/auth", m1.Then(http.HandlerFunc(c.Casbin.Add)))
	r.Del("/auth", m1.Then(http.HandlerFunc(c.Casbin.Delete)))
	r.Get("/auth/rule", m1.ThenFunc(c.Casbin.EditView(c)))
	r.Put("/auth", m1.Then(http.HandlerFunc(c.Casbin.Edit)))
	r.Get("/auth/roles4user/:userid", m1.ThenFunc(c.Casbin.AddRolesForUserView(c)))
	r.Post("/auth/roles4user/:userid", m1.Then(http.HandlerFunc(c.Casbin.AddRolesForUser)))

	r.Get("/files", m1.ThenFunc(c.File.Index(c)))
	r.Get("/zip", m1.ThenFunc(c.File.Zip(c)))
	r.Get("/dl", m1.ThenFunc(c.File.Download(c)))
	r.Get("/rebuild", m1.ThenFunc(c.File.Rebuild(c)))
	r.Get("/search", m1.ThenFunc(c.File.Search(c)))

	return m0.Then(r)
}
