package book

import (
	"word_of_wisdom_pow/internal/domain"
)

type Repository interface {
	GetRandomQuote() *domain.BookQuote
}
