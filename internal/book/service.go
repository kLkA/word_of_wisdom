package book

import (
	"word_of_wisdom_pow/internal/domain"
)

type Service interface {
	GetRandomQuote() *domain.BookQuote
}
