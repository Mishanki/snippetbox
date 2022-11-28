package main

import (
    "errors"
    "github.com/go-playground/form/v4"
    "net/http"
    "time"
)

func (app *application) newTemplateData(r *http.Request) *templateData {
    return &templateData{
        CurrentYear: time.Now().Year(),
        Flash: app.sessionManager.PopString(r.Context(), "flash"),
    }
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
    if err := r.ParseForm(); err != nil {
        return err
    }

    if err := app.formDecoder.Decode(dst, r.PostForm); err != nil {
        var invalidDecoderError *form.InvalidDecoderError
        if errors.As(err, &invalidDecoderError) {
            panic(err)
        }

        return err
    }

    return nil
}