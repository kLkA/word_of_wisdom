package service

import (
	"word_of_wisdom_pow/internal/book"
	"word_of_wisdom_pow/internal/domain"
)

type bookService struct {
	repo book.Repository
}

func NewBookService(repo book.Repository) book.Service {
	return &bookService{
		repo: repo,
	}
}
func (b bookService) GetRandomQuote() *domain.BookQuote {
	return b.repo.GetRandomQuote()
}
