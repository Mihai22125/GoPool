package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/pool/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createPoolForm))
	mux.Post("/pool/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createPool))
	mux.Get("/pool/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePoolForm)) 
	mux.Post("/pool/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePool))
	mux.Get("/pool/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showPool))
	mux.Get("/pool/update_options/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePoolOptionsForm))
	mux.Post("/pool/update_options/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePoolOptions))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))


	fileServer := http.FileServer(http.Dir(app.config.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}