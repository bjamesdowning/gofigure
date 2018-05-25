package controllers

import (
	"testing"

	"github.com/bjamesdowning/web2mongo/models"
	"github.com/stretchr/testify/assert"
)

func TestBookToJSON(t *testing.T) {
	book := models.Book{
		Title:  "Cloud Native",
		Author: "Writer",
		ISBN:   "1243532",
	}
	json := book.ToJSON()

	assert.Equal(t, `{"title":"Cloud Native","author":"Writer","isbn":"1243532"}`, string(json), "Book JSON Marshalling wrong")
}

func TestBookFromJSON(t *testing.T) {
	json := []byte(`{"title":"Cloud Native","author":"Writer","isbn":"1243532"}`)
	book := FromJSON(json)

	assert.Equal(t, models.Book{
		Title:  "Cloud Native",
		Author: "Writer",
		ISBN:   "1243532",
	}, book, "Book JSON Unmarshalling wrong")
}
