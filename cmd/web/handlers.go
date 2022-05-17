package main

import (
	"time"

	"github.com/mihai22125/goPool/pkg/models"
	"github.com/mihai22125/goPool/pkg/forms"
	"fmt"
	"strconv"
	"net/http"
)


func (app *application)home(w http.ResponseWriter, r *http.Request) {

	// userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	// if err != nil || userID < 1{
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	userID := 1

	pools, err := app.pools.GetAll(userID)
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Pools: pools,
	})
}

func (app *application) createPoolForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Pool: &models.Pool{},
		Form: forms.New(nil),
	})
}


func (app *application)createPool(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var startDate, endDate time.Time
	var nrOfOptions int64
	var singleVote bool

	form := forms.New(r.PostForm)
	form.Required("name", "nrOfOptions", "startDate", "endDate", "type")
	form.MaxLength("name", 100)
	form.IntegerRange("nrOfOptions", 1, 12)
	form.ValidDate("startDate", "2006-01-02")
	form.ValidDate("endDate", "2006-01-02")
	// TODO: rename to single-multiple
	form.PermittedValues("type", "0", "1")

	nrOfOptions, err = strconv.ParseInt(form.Get("nrOfOptions"), 0, 8)
	if err != nil {
		form.Errors.Add("nrOfOptions", "This field is invalid")
	}

	if startDate, err = time.Parse("2006-01-02", form.Get("startDate")); err != nil {
		form.Errors.Add("startDate", "This field is invalid")
	}

	if endDate, err = time.Parse("2006-01-02", form.Get("endDate")); err != nil {
		form.Errors.Add("endDate", "This field is invalid")
	}

	if form.Get("type") == "0" {
		singleVote = true
	}

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form, Pool: &models.Pool{}})
		return
	}
	
	poolConfig := models.PoolConfig{SingleVote: singleVote, StartDate: startDate, EndDate: endDate}

	pool := models.Pool{UserID: 1, Name: form.Get("name"), NumberOfOptions: int(nrOfOptions), PoolConfig: poolConfig}

	id, err := app.pools.Insert(pool)
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "flash", "Pool succesfully created!")

	// Redirect the user to the relevant page for the pool
	http.Redirect(w, r, fmt.Sprintf("/pool/update_options/%d", id), http.StatusSeeOther)
}


func (app *application) updatePoolForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}

	pool, err := app.pools.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "update.page.tmpl", &templateData{
		Form: forms.New(nil),
		Pool: pool,
	})
}

func (app *application)updatePool(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}

	pool, err := app.pools.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var startDate, endDate time.Time
	var nrOfOptions int64
	var singleVote bool

	form := forms.New(r.PostForm)
	form.Required("name", "nrOfOptions", "startDate", "endDate", "type")
	form.MaxLength("name", 100)
	form.IntegerRange("nrOfOptions", 1, 12)
	form.ValidDate("startDate", "2006-01-02")
	form.ValidDate("endDate", "2006-01-02")
	// TODO: rename to single-multiple
	form.PermittedValues("type", "0", "1")

	nrOfOptions, err = strconv.ParseInt(form.Get("nrOfOptions"), 0, 8)
	if err != nil {
		form.Errors.Add("nrOfOptions", "This field is invalid")
	}

	if startDate, err = time.Parse("2006-01-02", form.Get("startDate")); err != nil {
		form.Errors.Add("startDate", "This field is invalid")
	}

	if endDate, err = time.Parse("2006-01-02", form.Get("endDate")); err != nil {
		form.Errors.Add("endDate", "This field is invalid")
	}

	if form.Get("type") == "0" {
		singleVote = true
	}

	if !form.Valid() {
		app.render(w, r, "update.page.tmpl", &templateData{Form: form, Pool: pool})
		return
	}

	poolConfig := models.PoolConfig{PoolID: pool.ID, SingleVote: singleVote, StartDate: startDate, EndDate: endDate}
	pool.PoolConfig = poolConfig
	
	pool.Name = form.Get("name")
	pool.NumberOfOptions = int(nrOfOptions)

	_, err = app.pools.Update(pool)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// Redirect the user to the relevant page for the pool
	http.Redirect(w, r, fmt.Sprintf("/pool/%d", pool.ID), http.StatusSeeOther)
}

func (app *application) updatePoolOptionsForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}

	pool, err := app.pools.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "updateOptions.page.tmpl", &templateData{
		Form: forms.New(nil),
		Pool: pool,
	})
}


func (app *application)updatePoolOptions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	pool, err := app.pools.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for i := 0; i < pool.NumberOfOptions; i++ {
		form.Required(fmt.Sprintf("name[%d]", i), fmt.Sprintf("description[%d]", i))
		form.MaxLength(fmt.Sprintf("name[%d]", i), 100)
		form.MaxLength(fmt.Sprintf("description[%d]", i), 100)
	}

	if !form.Valid() {
		app.render(w, r, "updateOptions.page.tmpl", &templateData{Form: form, Pool: pool})
	}

	for i := 0; i < pool.NumberOfOptions; i++ {
		pool.PoolOptions[i].Option = form.Get(fmt.Sprintf("name[%d]", i))
		pool.PoolOptions[i].Description = form.Get(fmt.Sprintf("description[%d]", i))
	}

	_, err = app.pools.UpdateOptions(*pool)
	if err != nil {
		app.serveError(w, err)
		return
	}

	// Redirect the user to the relevant page for the pool
	http.Redirect(w, r, fmt.Sprintf("/pool/%d", id), http.StatusSeeOther)
}

func (app *application) showPool(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1{
		app.notFound(w)
		return
	}

	pool, err := app.pools.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Pool:  pool,
	})
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Add("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "userID", id)

	http.Redirect(w, r, "/pool/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
