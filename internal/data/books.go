package data

import (
	"Books/internal/validator"
	"time"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Authors   string    `json:"authors"`
	Rating    float64   `json:"rating"`
	ISBN      string    `json:"ISBN"`
	ISBN13    string    `json:"ISBN13"`
	Language  string    `json:"language,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Pages     Pages     `json:"pages,omitempty,string"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(book.Title != "", "title", "must be provided")
	v.Check(book.Authors != "", "authors", "must be provided")
	v.Check(len(book.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(book.ISBN != "", "ISBN", "must be provided")
	v.Check(book.ISBN13 != "", "ISBN13", "must be provided")
	v.Check(book.Rating != 0, "rating", "must be provided")
	v.Check(book.Rating > 0, "rating", "must be greater than 0")
	v.Check(book.Rating <= 5, "rating", "must be less than 5")
	v.Check(book.Pages != 0, "pages", "must be provided")
	v.Check(book.Pages > 0, "pages", "must be a positive integer")
	v.Check(book.Genres != nil, "genres", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
}
