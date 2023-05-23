package app

import (
	"word_of_wisdom_pow/internal/book"
	"word_of_wisdom_pow/internal/book/repository"
)

func (f *Factory) GetBookRepository() (*book.Repository, error) {
	if f.bookRepository == nil {
		repo, err := repository.NewBookRepository()
		if err != nil {
			return nil, err
		}
		f.bookRepository = &repo
	}
	return f.bookRepository, nil
}
