package app

import (
	"word_of_wisdom_pow/internal/book"
	_bookService "word_of_wisdom_pow/internal/book/service"
)

func (f *Factory) GetBookService() (book.Service, error) {
	if f.bookService == nil {
		bookRepo, err := f.GetBookRepository()
		if err != nil {
			return nil, err
		}
		bookService := _bookService.NewBookService(*bookRepo)
		f.bookService = bookService
	}
	return f.bookService, nil
}
