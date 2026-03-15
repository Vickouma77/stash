package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"stash.io/internal/models"
	"stash.io/internal/validator"
)

type StashCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type UserSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// home handler function with byte slice string
func (a *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := a.snippets.Latest()
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets

	a.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (a *Application) stashView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			a.ServerError(w, r, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.Snippet = snippet

	a.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (a *Application) stashCreate(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)

	data.Form = StashCreateForm{
		Expires: 365,
	}

	a.render(w, r, http.StatusOK, "create.tmpl.html", data)
}

func (a *Application) stashCreatePost(w http.ResponseWriter, r *http.Request) {
	var form StashCreateForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be emoty")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, 365")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form

		a.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "snippet created successfully")

	http.Redirect(w, r, fmt.Sprintf("/stash/view/%d", id), http.StatusSeeOther)
}

func (a *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = UserSignupForm{}

	a.render(w, r, http.StatusOK, "signup.tmpl.html", data)
}

func (a *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form UserSignupForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.ClientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form

		a.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	// Try to create a new user record in the database. If the email already
	// exists then add an error message to the form and re-display it.
	err = a.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "email address is already in use")

			data := a.newTemplateData(r)
			data.Form = form

			a.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			a.ServerError(w, r, err)
		}
		return
	}

	// Otherwise add a confirmation flash message to the session confirming that their signup worked.
	a.sessionManager.Put(r.Context(), "flash", "Your signup was successful, please login")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (a *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = UserLoginForm{}

	a.render(w, r, http.StatusOK, "login.tmpl.html", data)
}

func (a *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form UserLoginForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.ClientError(w, http.StatusBadRequest)
		return
	}
	// Validation Checking
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form

		a.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	// Checking valid credentials
	id, err := a.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorect")

			data := a.newTemplateData(r)
			data.Form = form

			a.render(w, r, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			a.ServerError(w, r, err)
		}
		return
	}

	err = a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/stash/create", http.StatusSeeOther)
}

func (a *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := a.sessionManager.RenewToken(r.Context())
	if err != nil {
		a.ServerError(w, r, err)
		return
	}

	a.sessionManager.Remove(r.Context(), "authenticatedUserID")

	a.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
