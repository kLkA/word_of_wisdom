package repository

import (
	_ "embed"
	"encoding/json"
	"math/rand"

	"word_of_wisdom_pow/internal/book"
	"word_of_wisdom_pow/internal/domain"
)

//go:embed quotes.json
var data []byte

type BookRepository struct {
	quotes []*domain.BookQuote
}

func NewBookRepository() (book.Repository, error) {
	var bookQuotes []*domain.BookQuote
	if err := json.Unmarshal(data, &bookQuotes); err != nil {
		return nil, err
	}
	r := &BookRepository{
		quotes: bookQuotes,
	}
	return r, nil
}

func (c BookRepository) GetRandomQuote() *domain.BookQuote {
	count := len(c.quotes)
	return c.quotes[rand.Intn(count)]
}
