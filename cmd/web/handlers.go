package main

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "runtime/debug"
    "snippetbox/internal/models"
    "snippetbox/internal/validator"
    "strconv"
)

type snippetCreateForm struct {
    Title string
    Content string
    Expires int
    validator.Validator `form:"-"`
}

type userSignupForm struct {
    Name string `form:"name"`
    Email string `form:"email"`
    Password string `form:"password"`
    validator.Validator `form:"-"`
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

    data.Form = snippetCreateForm{
        Expires: 365,
    }

    app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

    var form snippetCreateForm
    if err := app.decodePostForm(r, &form); err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    fmt.Println(form)

    form.CheckField(validator.NotBlank(form.Title), "title", "This field can`t be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field can`t be more than 100 characters long")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
    form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

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

func (app *application) userSignup (w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = userSignupForm{}
    app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *application) userSignupPost (w http.ResponseWriter, r *http.Request) {
    var form userSignupForm

    if err := app.decodePostForm(r, &form); err != nil {
        app.clientError(w, http.StatusBadRequest)
        return
    }

    form.CheckField(validator.NotBlank(form.Name), "name", "This field can`t be blank")
    form.CheckField(validator.NotBlank(form.Email), "email", "This field can`t be blank")
    form.CheckField(validator.NotBlank(form.Password), "password", "This field can`t be blank")
    form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
    form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
        return
    }

    fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin (w http.ResponseWriter, r *http.Request) {
    fmt.Println(w, "Display a HTML form for logging in a user...")
}

func (app *application) userLoginPost (w http.ResponseWriter, r *http.Request) {
    fmt.Println(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost (w http.ResponseWriter, r *http.Request) {
    fmt.Println(w, "Logout the user...")
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
    ts, ok := app.templateCache[page]
    if !ok {
        err := fmt.Errorf("the template %s does not exists", page)
        app.serverError(w, err)
        return
    }

    buf := new(bytes.Buffer)

    err := ts.ExecuteTemplate(buf, "base", data)
    if err != nil {
        app.serverError(w, err)
        return
    }

    w.WriteHeader(status)
    buf.WriteTo(w)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    app.errorLog.Output(2, trace)

    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
    http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
    app.clientError(w, http.StatusNotFound)
}