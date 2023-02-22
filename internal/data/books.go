package data

import "time"

type Books struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Authors         string    `json:"authors"`
	Rating          float64   `json:"rating"`
	ISBN            string    `json:"ISBN"`
	ISBN13          string    `json:"ISBN13"`
	Language        string    `json:"language,omitempty"`
	Genres          []string  `json:"genres,omitempty"`
	Pages           Pages     `json:"pages,omitempty,string"`
	PublicationDate time.Time `json:"-"`
}
