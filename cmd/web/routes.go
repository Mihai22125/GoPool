package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/pool/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createPoolForm))
	mux.Post("/pool/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createPool))
	mux.Get("/pool/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePoolForm)) 
	mux.Post("/pool/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePool))
	mux.Get("/pool/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showPool))
	mux.Get("/pool/update_options/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePoolOptionsForm))
	mux.Post("/pool/update_options/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.updatePoolOptions))

	mux.Get("/pool/results/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showPoolResults))


	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	// admin routes
	mux.Get("/machine/create", dynamicMiddleware.Append(app.requireAuthenticatedUser, app.adminUser).ThenFunc(app.createMachineForm))
	mux.Post("/machine/create", dynamicMiddleware.Append(app.requireAuthenticatedUser, app.adminUser).ThenFunc(app.createMachine))
	mux.Get("/machine/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser, app.adminUser).ThenFunc(app.updateMachineForm))
	mux.Post("/machine/update/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser, app.adminUser).ThenFunc(app.updateMachine))
	mux.Post("/machine/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser, app.adminUser).ThenFunc(app.deleteMachine))
	
	mux.Get("/machine/", dynamicMiddleware.Append(app.requireAuthenticatedUser, app.adminUser).ThenFunc(app.showMachines))
	mux.Get("/machine/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser, app.adminUser).ThenFunc(app.showMachine))

	mux.Post("/vote", http.HandlerFunc(app.createVote))


	fileServer := http.FileServer(http.Dir(app.config.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}

func (app *application) restRoutes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Post("/vote", http.HandlerFunc(app.createVote))

	return standardMiddleware.Then(mux)
}