package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.brainwhat/internal/models"
	"snippetbox.brainwhat/internal/validator"
)

// We include these tags to show the decoder how to map values from the form
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // This one means innore the field
}

// Small func for testing
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// We preset the Expires radio field so data isn't empty and
	// {{ with .Form.FieldErrors.title }} in create.tmpl cause 500 err
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.userError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignUpForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.userError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MatchesRX(form.Email, validator.EmailRX), "email", "Invalid email")
	form.CheckField(validator.MinChars(form.Name, 3), "name", "Name too short")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password too short")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = app.users.Create(form.Email, form.Name, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldErrors("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "New user created")

	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
}

type UserSignInForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) userSignIn(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserSignInForm{}
	app.render(w, http.StatusOK, "signin.tmpl", data)
}

func (app *application) userSignInPost(w http.ResponseWriter, r *http.Request) {
	var form UserSignInForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.userError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.MatchesRX(form.Email, validator.EmailRX), "email", "Invalid email")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signin.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signin.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	url := app.sessionManager.PopString(r.Context(), "requestedURL")
	if url != "" {
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogOutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	u, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "user/signin", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
	}

	data := app.newTemplateData(r)
	data.User = u

	app.render(w, http.StatusOK, "account.tmpl", data)

}

type accountPasswordUpdateForm struct {
	CurPass             string `form:"curPass"`
	NewPass             string `form:"newPass"`
	NewPassConfirm      string `form:"newPassConfirm"`
	validator.Validator `form:"-"`
}

func (app *application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = accountPasswordUpdateForm{}
	app.render(w, http.StatusOK, "passwordChange.tmpl", data)
}

func (app *application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	var form accountPasswordUpdateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.userError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.CurPass), "curPass", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPass), "newPass", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.NewPassConfirm), "nePassConfirm", "This field cannot be blank")

	form.CheckField(validator.MinChars(form.NewPass, 8), "newPass", "New password must be 8+ characters long")

	form.CheckField(validator.MatchesString(form.NewPass, form.NewPassConfirm), "newPass", "Passwords don't match")
	form.CheckField(validator.MatchesString(form.NewPass, form.NewPassConfirm), "newPassConfirm", "Passwords don't match")

	if !form.Valid() {
		app.infoLog.Print("Validator erro")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "passwordChange.tmpl", data)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.UpdatePassword(id, form.CurPass, form.NewPass)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Password updated successfully")

	http.Redirect(w, r, "/user/account", http.StatusSeeOther)
}
