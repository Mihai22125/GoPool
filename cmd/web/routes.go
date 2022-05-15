package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/pool/create", dynamicMiddleware.ThenFunc(app.createPoolForm))
	mux.Post("/pool/create", dynamicMiddleware.ThenFunc(app.createPool))
	mux.Get("/pool/update/:id", dynamicMiddleware.ThenFunc(app.updatePoolForm)) //TODO: replace handler function
	mux.Post("/pool/update/:id", dynamicMiddleware.ThenFunc(app.updatePool))    //TODO: replace handler function
	mux.Get("/pool/:id", dynamicMiddleware.ThenFunc(app.showPool))
	mux.Get("/pool/update_options/:id", dynamicMiddleware.ThenFunc(app.updatePoolOptionsForm))
	mux.Post("/pool/update_options/:id", dynamicMiddleware.ThenFunc(app.updatePoolOptions))

	fileServer := http.FileServer(http.Dir(app.config.StaticDir))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}