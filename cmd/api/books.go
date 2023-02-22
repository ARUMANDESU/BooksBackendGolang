package main

import (
	"Books/internal/data"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	books := data.Books{
		ID:              id,
		Title:           "Harry Potter",
		Authors:         "J.K. Rowling",
		Rating:          4.57,
		ISBN:            "0439785960",
		ISBN13:          "9780439785969",
		Language:        "",
		Pages:           652,
		PublicationDate: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
