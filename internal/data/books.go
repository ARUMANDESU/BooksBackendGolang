package data

import (
	"Books/internal/validator"
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
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
	v.Check(book.Language != "", "language", "must be provided")
	v.Check(len(book.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(book.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(book.Genres), "genres", "must not contain duplicate values")
}

type BookModel struct {
	DB *sql.DB
}

func (b BookModel) Insert(book *Book) error {
	query := `
			INSERT INTO books (title, authors,rating, pages, genres,isbn,isbn13,language)
			VALUES ($1, $2, $3, $4, $5,$6,$7,$8)
			RETURNING id, created_at, version`

	args := []any{book.Title, book.Authors, book.Rating, book.Pages, pq.Array(book.Genres), book.ISBN, book.ISBN13, book.Language}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return b.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
			SELECT  id, created_at, title, authors,rating, pages, genres,isbn,isbn13,language,version
			FROM books
			WHERE id = $1`

	var book Book
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	err := b.DB.QueryRowContext(ctx, query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Authors,
		&book.Rating,
		&book.Pages,
		pq.Array(&book.Genres),
		&book.ISBN,
		&book.ISBN13,
		&book.Language,
		&book.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &book, nil
}

func (b BookModel) Update(book *Book) error {
	query := `
			UPDATE books
			SET title = $1, authors = $2, pages = $3, rating=$4, genres = $5, isbn=$6, isbn13=$7, language=$8, version = version + 1
			WHERE id = $9 and version = $10
			RETURNING version`

	args := []any{
		book.Title,
		book.Authors,
		book.Pages,
		book.Rating,
		pq.Array(book.Genres),
		book.ISBN,
		book.ISBN13,
		book.Language,
		book.ID,
		book.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := b.DB.QueryRowContext(ctx, query, args...).Scan(&book.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err

		}
	}
	return nil
}

func (b BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
			DELETE FROM books
			WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := b.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
