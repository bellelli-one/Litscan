package repository

import (
	"RIP/internal/app/ds"
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GET /api/books - список книг с фильтрацией
func (r *Repository) BooksList(title string) ([]ds.Books, int64, error) {
	var factors []ds.Books
	var total int64

	query := r.db.Model(&ds.Books{})
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	factorsQuery := query.Order("id asc")
	if err := factorsQuery.Find(&factors).Error; err != nil {
		return nil, 0, err
	}

	return factors, total, nil
}

// GET /api/books/:id - одна книга
func (r *Repository) GetBookByID(id int) (*ds.Books, error) {
	var book ds.Books
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// POST /api/books - создание книги
func (r *Repository) CreateBook(book *ds.Books) error {
	return r.db.Create(book).Error
}

// PUT /api/books/:id - обновление книги
func (r *Repository) UpdateBook(id uint, req ds.BookUpdateRequest) (*ds.Books, error) {
	var book ds.Books
	if err := r.db.First(&book, id).Error; err != nil {
		return nil, err
	}

	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Text != nil {
		book.Text = *req.Text
	}
	if req.AvgWordLen != nil {
		book.AvgWordLen = *req.AvgWordLen
	}
	if req.LexicalDiversity != nil {
		book.LexicalDiversity = *req.LexicalDiversity
	}
	if req.ConjunctionFreq != nil {
		book.ConjunctionFreq = *req.ConjunctionFreq
	}
	if req.AvgSentenceLen != nil {
		book.AvgSentenceLen = *req.AvgSentenceLen
	}

	if err := r.db.Save(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// DELETE /api/books/:id - удаление книги
func (r *Repository) DeleteBook(id uint) error {
	var book ds.Books
	var imageURLToDelete string

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&book, id).Error; err != nil {
			return err
		}
		if book.Image != nil {
			imageURLToDelete = *book.Image
		}
		if err := tx.Delete(&ds.Books{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	if imageURLToDelete != "" {
		parsedURL, err := url.Parse(imageURLToDelete)
		if err != nil {
			log.Printf("ERROR: could not parse image URL for deletion: %v", err)
			return nil
		}

		objectName := strings.TrimPrefix(parsedURL.Path, fmt.Sprintf("/%s/", r.bucketName))

		err = r.minioClient.RemoveObject(context.Background(), r.bucketName, objectName, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("ERROR: failed to delete object '%s' from MinIO: %v", objectName, err)
		}
	}

	return nil
}

// POST /api/analysebooks/draft/books/:book_id - добавление книги в черновик
func (r *Repository) AddBookToDraft(userID, bookID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var appl ds.AnalyseBooks
		err := tx.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&appl).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newAppl := ds.AnalyseBooks{
					CreatorID:    userID,
					Status:       ds.StatusDraft,
					CreationDate: time.Now(),
				}
				if err := tx.Create(&newAppl).Error; err != nil {
					return fmt.Errorf("failed to create draft appl: %w", err)
				}
				appl = newAppl
			} else {
				return err
			}
		}

		var count int64
		tx.Model(&ds.BookToAppl{}).Where("appl_id = ? AND book_id = ?", appl.ID, bookID).Count(&count)
		if count > 0 {
			return errors.New("book already in appl")
		}

		link := ds.BookToAppl{
			ApplID: appl.ID,
			BookID: bookID,
		}

		if err := tx.Create(&link).Error; err != nil {
			return fmt.Errorf("failed to add book to appl: %w", err)
		}

		if err := tx.Model(&ds.Books{}).Where("id = ?", bookID).Update("status", true).Error; err != nil {
			return fmt.Errorf("failed to update book status: %w", err)
		}
		return nil
	})
}

// POST /api/books/:id/image - загрузка изображения фактора
func (r *Repository) UploadBookImage(bookID uint, fileHeader *multipart.FileHeader) (string, error) {
	var finalImageURL string
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var book ds.Books
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, bookID).Error; err != nil {
			return fmt.Errorf("book with id %d not found: %w", bookID, err)
		}

		const imagePathPrefix = "Images/"

		if book.Image != nil && *book.Image != "" {
			oldImageURL, err := url.Parse(*book.Image)
			if err == nil {
				oldObjectName := strings.TrimPrefix(oldImageURL.Path, fmt.Sprintf("/%s/", r.bucketName))
				r.minioClient.RemoveObject(context.Background(), r.bucketName, oldObjectName, minio.RemoveObjectOptions{})
			}
		}

		fileName := filepath.Base(fileHeader.Filename)
		objectName := imagePathPrefix + fileName

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = r.minioClient.PutObject(context.Background(), r.bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})

		if err != nil {
			return fmt.Errorf("failed to upload to minio: %w", err)
		}

		imageURL := fmt.Sprintf("http://%s/%s/%s", r.minioEndpoint, r.bucketName, objectName)

		if err := tx.Model(&book).Update("image", imageURL).Error; err != nil {
			return fmt.Errorf("failed to update book image url in db: %w", err)
		}

		finalImageURL = imageURL
		return nil
	})
	if err != nil {
		return "", err
	}
	return finalImageURL, nil
}

// func (r *Repository) GetAllBooks() ([]ds.Books, error) {
// 	var books []ds.Books

// 	err := r.db.Find(&books).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(books) == 0 {
// 		return nil, fmt.Errorf("books not found")
// 	}
// 	return books, nil
// }

// func (r *Repository) SearchBooksByName(title string) ([]ds.Books, error) {
// 	var books []ds.Books
// 	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&books).Error // добавили условие
// 	if err != nil {
// 		return nil, err
// 	}
// 	return books, nil
// }

// func (r *Repository) GetBookByID(id int) (*ds.Books, error) {
// 	var books ds.Books
// 	err := r.db.First(&books, id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &books, nil
// }
