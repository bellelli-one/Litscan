package repository

import (
	"RIP/internal/app/ds"
	"fmt"
)

func (r *Repository) GetAllBooks() ([]ds.Books, error) {
	var books []ds.Books

	err := r.db.Find(&books).Error
	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, fmt.Errorf("books not found")
	}
	return books, nil
}

func (r *Repository) SearchBooksByName(title string) ([]ds.Books, error) {
	var books []ds.Books
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&books).Error // добавили условие
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (r *Repository) GetBookByID(id int) (*ds.Books, error) {
	var books ds.Books
	err := r.db.First(&books, id).Error
	if err != nil {
		return nil, err
	}
	return &books, nil
}