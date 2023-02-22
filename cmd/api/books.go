package main

import (
	"Books/internal/data"
	"Books/internal/validator"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string     `json:"title"`
		Authors  string     `json:"authors"`
		ISBN     string     `json:"ISBN"`
		ISBN13   string     `json:"ISBN13"`
		Language string     `json:"language"`
		Genres   []string   `json:"genres"`
		Rating   float64    `json:"rating"`
		Pages    data.Pages `json:"pages"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	book := &data.Book{
		Title:    input.Title,
		Authors:  input.Authors,
		Rating:   input.Rating,
		ISBN:     input.ISBN,
		ISBN13:   input.ISBN13,
		Language: input.Language,
		Genres:   input.Genres,
		Pages:    input.Pages,
	}

	v := validator.New()

	if data.ValidateBook(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	books := data.Book{
		ID:        id,
		Title:     "Harry Potter",
		Authors:   "J.K. Rowling",
		Rating:    4.57,
		ISBN:      "0439785960",
		ISBN13:    "9780439785969",
		Language:  "",
		Pages:     652,
		CreatedAt: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
