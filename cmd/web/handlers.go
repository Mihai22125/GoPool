package main

import (
	"encoding/json"
	"net"
	"sort"
	"time"

	"fmt"
	"net/http"
	"strconv"

	"github.com/mihai22125/goPool/pkg/forms"
	"github.com/mihai22125/goPool/pkg/models"
)

//TODO: add default page for not authenticated users
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	user := app.authenticatedUser(r)
	var userID int

	if user != nil {
		userID = user.ID
	} else {
		userID = 1
	}

	pools, err := app.pools.GetAll(userID)
	if err != nil {
		app.errorLog.Println(err)
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

func (app *application) createPool(w http.ResponseWriter, r *http.Request) {
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
	form.ValidDate("startDate", "2006-01-02T15:04")
	form.ValidDate("endDate", "2006-01-02T15:04")
	// TODO: rename to single-multiple
	form.PermittedValues("type", "0", "1")

	nrOfOptions, err = strconv.ParseInt(form.Get("nrOfOptions"), 0, 8)
	if err != nil {
		form.Errors.Add("nrOfOptions", "This field is invalid")
	}

	fmt.Println("startDateString: ", form.Get("startDate"))
	if startDate, err = time.Parse("2006-01-02T15:04", form.Get("startDate")); err != nil {
		form.Errors.Add("startDate", "This field is invalid")
	}
	fmt.Println("startDate: ", startDate)

	if time.Now().After(startDate) {
		form.Errors.Add("startDate", "Start Date must be in future")
	}

	if endDate, err = time.Parse("2006-01-02T15:04", form.Get("endDate")); err != nil {
		app.errorLog.Println(err)
		form.Errors.Add("endDate", "This field is invalid")
	}

	if time.Now().After(endDate) {
		form.Errors.Add("endDate", "End Date must be in future")
	}

	hours := endDate.Sub(startDate).Hours()
	if hours < 0 {
		form.Errors.Add("endDate", "endDate must be after startDate")
	}

	if hours > 2 {
		form.Errors.Add("endDate", "duration must be less than 2 hours")
	}

	if form.Get("type") == "0" {
		singleVote = true
	}

	machineID, err := app.machines.GetNextAvailable(startDate, endDate)
	if err != nil {
		app.serveError(w, err)
		return
	}

	if machineID == 0 {
		form.Errors.Add("startDate", "this period is unavailable")
	}

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form, Pool: &models.Pool{}})
		return
	}

	poolConfig := models.PoolConfig{SingleVote: singleVote, StartDate: startDate, EndDate: endDate}

	user := app.authenticatedUser(r)

	pool := models.Pool{UserID: user.ID, Name: form.Get("name"), NumberOfOptions: int(nrOfOptions), PoolConfig: poolConfig}

	//TODO: make transaction at inserting pool - session
	id, err := app.pools.Insert(pool)
	if err != nil {
		app.serveError(w, err)
		return
	}

	_, err = app.sessions.Insert(id, machineID)
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
	if err != nil || id < 1 {
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

	user := app.authenticatedUser(r)
	if pool.UserID != user.ID {
		app.notFound(w)
		return
	}

	app.render(w, r, "update.page.tmpl", &templateData{
		Form: forms.New(nil),
		Pool: pool,
	})
}

func (app *application) updatePool(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
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

	if time.Now().After(pool.PoolConfig.StartDate) {
		return
	}

	user := app.authenticatedUser(r)
	if pool.UserID != user.ID {
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var startDate, endDate time.Time
	var singleVote bool

	form := forms.New(r.PostForm)
	form.Required("name", "startDate", "endDate", "type")
	form.MaxLength("name", 100)
	form.ValidDate("startDate", "2006-01-02")
	form.ValidDate("endDate", "2006-01-02")
	// TODO: rename to single-multiple
	form.PermittedValues("type", "0", "1")

	if startDate, err = time.Parse("2006-01-02", form.Get("startDate")); err != nil {
		form.Errors.Add("startDate", "This field is invalid")
	}

	if endDate, err = time.Parse("2006-01-02T15:04", form.Get("endDate")); err != nil {
		app.errorLog.Println(err)
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
	if err != nil || id < 1 {
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

	user := app.authenticatedUser(r)
	if pool.UserID != user.ID {
		app.notFound(w)
		return
	}

	app.render(w, r, "updateOptions.page.tmpl", &templateData{
		Form: forms.New(nil),
		Pool: pool,
	})
}

func (app *application) updatePoolOptions(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
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

	user := app.authenticatedUser(r)
	if pool.UserID != user.ID {
		app.notFound(w)
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
	if err != nil || id < 1 {
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

	user := app.authenticatedUser(r)
	if pool.UserID != user.ID {
		app.notFound(w)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Pool: pool,
	})
}

func (app *application) showPoolResults(w http.ResponseWriter, r *http.Request) {
	var results []*models.Result

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
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

	var total int

	for _, option := range pool.PoolOptions {
		count, err := app.votes.CountByOptionID(pool.ID, option.ID)
		if err != nil {
			app.serveError(w, err)
			return
		}
		results = append(results, &models.Result{Option: option, Count: count})
		total += count
	}

	for i := range results {
		if results[i].Count == 0 {
			results[i].Percentage = 0
		} else {
			results[i].Percentage = float64(results[i].Count) / float64(total) * 100.0
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Count > results[j].Count
	})

	app.render(w, r, "show_pool_results.page.tmpl", &templateData{
		Pool:    pool,
		Results: results,
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

func (app *application) createVote(w http.ResponseWriter, r *http.Request) {
	var voteRequest models.VoteRequest

	err := json.NewDecoder(r.Body).Decode(&voteRequest)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Printf("%+v\n", voteRequest)

	machine, err := app.machines.Get(voteRequest.MachineID)
	if err != nil {
		app.serveError(w, err)
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Printf("machine: %+v\n", machine)

	clientHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		app.serveError(w, err)
		app.errorLog.Println(err)
		return
	}
	// if clientHost != machine.IPAdrres {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	app.errorLog.Println(err)
	// 	return
	// }

	app.infoLog.Printf("clienthost: %v", clientHost)

	session, err := app.sessions.GetCurrentForMachine(machine.ID)
	if err == models.ErrNoRecord {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Println(err)
		return
	} else if err != nil {
		app.serveError(w, err)
		app.errorLog.Println(err)
		return
	}
	app.infoLog.Printf("session: %+v\n", session)

	optionID, err := app.pools.GetOptionID(session.PoolID, voteRequest.Text)
	app.infoLog.Printf("option id: %v\n", optionID)
	if err == models.ErrNoRecord {
		app.clientError(w, http.StatusBadRequest)
		app.errorLog.Println(err)
		return
	} else if err != nil {
		app.serveError(w, err)
		app.errorLog.Println(err)
		return
	}

	_, err = app.votes.Insert(session.PoolID, optionID, machine.ID, voteRequest.From)
	if err != nil {
		app.serveError(w, err)
		app.errorLog.Println(err)
		return
	}
}

func (app *application) createMachineForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create_machine.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createMachine(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("ipAddress", "phoneNumber")
	form.MaxLength("ipAddress", 15)
	form.MaxLength("phoneNumber", 15)

	if !form.Valid() {
		app.render(w, r, "create_machine.page.tmpl", &templateData{Form: form})
		return
	}

	machine := models.Machine{IPAdrres: form.Get("ipAddress"), PhoneNumber: form.Get("phoneNumber")}

	//TODO: make transaction at inserting pool - session
	_, err = app.machines.Insert(machine)
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "flash", "Machine succesfully added!")

	// Redirect the user to the relevant page for the pool
	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

func (app *application) showMachine(w http.ResponseWriter, r *http.Request) {
}

func (app *application) showMachines(w http.ResponseWriter, r *http.Request) {
	machines, err := app.machines.GetAll()
	if err != nil {
		app.errorLog.Println(err)
		app.serveError(w, err)
		return
	}

	app.render(w, r, "show_machines.page.tmpl", &templateData{
		Machines: machines,
	})
}

func (app *application) updateMachineForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	machine, err := app.machines.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.render(w, r, "update_machine.page.tmpl", &templateData{
		Form:    forms.New(nil),
		Machine: machine,
	})
}

func (app *application) updateMachine(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	machine, err := app.machines.Get(id)
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

	form := forms.New(r.PostForm)
	form.Required("ipAddress", "phoneNumber")
	form.MaxLength("ipAddress", 15)
	form.MaxLength("phoneNumber", 15)

	if !form.Valid() {
		app.render(w, r, "update_machine.page.tmpl", &templateData{Form: form, Machine: machine})
		return
	}

	machine.IPAdrres = form.Get("ipAddress")
	machine.PhoneNumber = form.Get("phoneNumber")

	_, err = app.machines.Update(machine)
	if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "flash", "Machine succesfully updated!")

	// Redirect the user to the relevant page for the pool
	http.Redirect(w, r, fmt.Sprintf("/machine"), http.StatusSeeOther)
}

func (app *application) deleteMachine(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	_, err = app.machines.Delete(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveError(w, err)
		return
	}

	app.session.Put(r, "flash", "Machine succesfully deleted!")

	http.Redirect(w, r, fmt.Sprintf("/machine"), http.StatusSeeOther)
}
