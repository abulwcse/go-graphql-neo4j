// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type NewAuthor struct {
	ID          int     `json:"id"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	DateOfBirth *string `json:"dateOfBirth"`
}

type NewBook struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Isbn     string `json:"isbn"`
	Language string `json:"language"`
	AuthorID int    `json:"authorId"`
}
