package repository

import (
	"RIP/internal/app/ds"
	"errors"
)

func (r *Repository) GetDraftAppl(userID uint) (*ds.Application, error) {
	var appl ds.Application

	err := r.db.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&appl).Error
	if err != nil {
		return nil, err
	}
	return &appl, nil
}

func (r *Repository) CreateAppl(appl *ds.Application) error {
	return r.db.Create(appl).Error
}

func (r *Repository) AddBookToAppl(applID, bookID uint) error {
	var count int64

	r.db.Model(&ds.BookToAppl{}).Where("appl_id = ? AND book_id = ?", applID, bookID).Count(&count)
	if count > 0 {
		return errors.New("book already in application")
	}

	link := ds.BookToAppl{
		ApplID: applID,
		BookID: bookID,
	}
	return r.db.Create(&link).Error
}

func (r *Repository) GetApplWithBooks(applID uint) (*ds.Application, error) {
	var appl ds.Application

	err := r.db.Preload("BooksLink.Book").First(&appl, applID).Error
	if err != nil {
		return nil, err
	}

	if appl.Status == ds.StatusDeleted {
		return nil, errors.New("book page not found or has been deleted")
	}

	return &appl, nil
}

func (r *Repository) LogicallyDeleteAppl(applID uint) error {
	result := r.db.Exec("UPDATE applications SET status = ? WHERE id = ?", ds.StatusDeleted, applID)
	return result.Error
}
