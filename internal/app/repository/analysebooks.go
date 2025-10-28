package repository

import (
	"RIP/internal/app/ds"
	"errors"
	"math"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// GET /api/analysebooks/cart - иконка корзины
func (r *Repository) GetDraftAnalyseBooks(userID uint) (*ds.AnalyseBooks, error) {
	var appl ds.AnalyseBooks
	err := r.db.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&appl).Error
	if err != nil {
		return nil, err
	}
	return &appl, nil
}

// GET /api/analysebooks/cart - иконка корзины
// GET /api/analysebooks/:id - одна заявка с услугами
func (r *Repository) GetAnalyseBooksWithBooks(applID uint) (*ds.AnalyseBooks, error) {
	var appl ds.AnalyseBooks
	err := r.db.Preload("BooksLink.Book").Preload("Creator").Preload("Moderator").First(&appl, applID).Error
	if err != nil {
		return nil, err
	}

	if appl.Status == ds.StatusDeleted {
		return nil, errors.New("appl page not found or has been deleted")
	}

	return &appl, nil
}

// GET /api/analysebooks - список заявок с фильтрацией
func (r *Repository) AnalyseBooksListFiltered(userID uint, isModerator bool, status, from, to string) ([]ds.AnalyseBooksDTO, error) {
	var appList []ds.AnalyseBooks
	query := r.db.Preload("Creator").Preload("Moderator")

	query = query.Where("status != ? AND status != ?", ds.StatusDeleted, ds.StatusDraft)

	if !isModerator {
		query = query.Where("creator_id = ?", userID)
	}

	if status != "" {
		if statusInt, err := strconv.Atoi(status); err == nil {
			query = query.Where("status = ?", statusInt)
		}
	}

	if from != "" {
		if fromTime, err := time.Parse("2006-01-02", from); err == nil {
			query = query.Where("forming_date >= ?", fromTime)
		}
	}

	if to != "" {
		if toTime, err := time.Parse("2006-01-02", to); err == nil {
			query = query.Where("forming_date <= ?", toTime)
		}
	}

	if err := query.Find(&appList).Error; err != nil {
		return nil, err
	}

	var result []ds.AnalyseBooksDTO
	for _, app := range appList {
		dto := ds.AnalyseBooksDTO{
			ID:               app.ID,
			Status:           app.Status,
			CreationDate:     app.CreationDate,
			CreatorID:        app.Creator.ID,
			ModeratorID:      nil,
			FormingDate:      app.FormingDate,
			CompletionDate:   app.ComplitionDate,
			AwgWordLen:       app.AwgWordLen,
			LexicalDiversity: app.LexicalDiversity,
			ConjunctionFreq:  app.ConjunctionFreq,
			AvgSentenceLen:   app.AvgSentenceLen,
			Response:         app.Response,
		}

		if app.ModeratorID != nil {
			dto.ModeratorID = &app.Moderator.ID
		}
		result = append(result, dto)
	}
	return result, nil
}

// PUT /api/analysebooks/:id - изменение полей заявки
func (r *Repository) UpdateAnalyseBooksUserFields(id uint, req ds.AnalyseBooksUpdateRequest) error {
	updates := make(map[string]interface{})

	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Moderator != nil {
		updates["moderator"] = *req.Moderator
	}
	if req.FormingDate != nil {
		updates["forming_date"] = *req.FormingDate
	}
	if req.CompletionDate != nil {
		updates["complition_date"] = *req.CompletionDate
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&ds.AnalyseBooks{}).Where("id = ?", id).Updates(updates).Error
}

// PUT /api/analysebooks/:id/form - сформировать заявку
func (r *Repository) FormAnalyseBooks(id uint, creatorID uint) error {
	var app ds.AnalyseBooks
	if err := r.db.First(&app, id).Error; err != nil {
		return err
	}

	if app.CreatorID != creatorID {
		return errors.New("only creator can form application")
	}

	if app.Status != ds.StatusDraft {
		return errors.New("only draft application can be formed")
	}

	// Проверяем, есть ли связанные книги
	var booksCount int64
	if err := r.db.Model(&ds.BookToAppl{}).Where("appl_id = ?", id).Count(&booksCount).Error; err != nil {
		return err
	}

	if booksCount == 0 {
		return errors.New("at least one book is required to form application")
	}

	now := time.Now()
	return r.db.Model(&app).Updates(map[string]interface{}{
		"status":       ds.StatusFormed,
		"forming_date": now,
	}).Error
}

// PUT /api/analyse-books/:id/resolve - завершить/отклонить заявку
func (r *Repository) ResolveAnalyseBooks(id uint, moderatorID uint, action string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		var app ds.AnalyseBooks
		if err := tx.Preload("BooksLink.Book").First(&app, id).Error; err != nil {
			return err
		}

		if app.Status != ds.StatusFormed {
			return errors.New("only formed application can be resolved")
		}

		now := time.Now()
		updates := map[string]interface{}{
			"moderator_id":    moderatorID,
			"complition_date": now,
		}

		switch action {
		case "complete":
			{
				updates["status"] = ds.StatusCompleted
				// Здесь можно добавить расчеты для анализа книг, если нужно
				// awgWordLen, lexicalDiversity, etc.
			}
		case "reject":
			{
				updates["status"] = ds.StatusRejected
			}
		default:
			{
				return errors.New("invalid action, must be 'complete' or 'reject'")
			}
		}

		if err := tx.Model(&app).Updates(updates).Error; err != nil {
			return err
		}

		var bookIDs []uint
		for _, link := range app.BooksLink {
			bookIDs = append(bookIDs, link.BookID)
		}

		if len(bookIDs) > 0 {
			if err := tx.Model(&ds.Books{}).Where("id IN ?", bookIDs).Update("status", false).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Функция расчета схожести по расстоянию Канберра
func (r *Repository) calculateBookSimilarity(app ds.AnalyseBooks, book ds.Books) float64 {
	// Вектор из заявкиt
	appVector := []float64{
		app.AwgWordLen,
		app.LexicalDiversity,
		app.ConjunctionFreq,
		app.AvgSentenceLen,
	}

	// Вектор из книги
	bookVector := []float64{
		book.AvgWordLen,
		book.LexicalDiversity,
		book.ConjunctionFreq,
		book.AvgSentenceLen,
	}

	// Расчет расстояния Канберра
	distance := 0.0
	for i := 0; i < len(appVector); i++ {
		numerator := math.Abs(appVector[i] - bookVector[i])
		denominator := math.Abs(appVector[i]) + math.Abs(bookVector[i])

		// Избегаем деления на ноль
		if denominator != 0 {
			distance += numerator / denominator
		}
	}

	// Преобразование расстояния в вероятность схожести (0-100%)
	// Чем меньше расстояние, тем больше вероятность
	similarity := (1 - (distance / float64(len(appVector)))) * 100

	// Ограничиваем значения от 0 до 100
	if similarity < 0 {
		similarity = 0
	}
	if similarity > 100 {
		similarity = 100
	}

	return similarity
}

// DELETE /api/analysebooks/:id - удаление заявки
func (r *Repository) LogicallyDeleteAnalyseBooks(appID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var app ds.AnalyseBooks

		if err := tx.Preload("BooksLink").First(&app, appID).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"status":       ds.StatusDeleted,
			"forming_date": time.Now(),
		}

		if err := tx.Model(&ds.AnalyseBooks{}).Where("id = ?", appID).Updates(updates).Error; err != nil {
			return err
		}

		var bookIDs []uint
		for _, link := range app.BooksLink {
			bookIDs = append(bookIDs, link.BookID)
		}

		if len(bookIDs) > 0 {
			if err := tx.Model(&ds.Books{}).Where("id IN ?", bookIDs).Update("status", false).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DELETE /api/analysebooks/:id/books/:book_id - удаление книги из заявки
func (r *Repository) RemoveBookFromAnalyseBooks(appID, bookID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		result := tx.Where("appl_id = ? AND book_id = ?", appID, bookID).Delete(&ds.BookToAppl{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("book not found in this application")
		}

		if err := tx.Model(&ds.Books{}).Where("id = ?", bookID).Update("status", false).Error; err != nil {
			return err
		}

		var remainingCount int64
		if err := tx.Model(&ds.BookToAppl{}).Where("appl_id = ?", appID).Count(&remainingCount).Error; err != nil {
			return err
		}

		if remainingCount == 0 {
			updates := map[string]interface{}{
				"status":       ds.StatusDeleted,
				"forming_date": time.Now(),
			}
			if err := tx.Model(&ds.AnalyseBooks{}).Where("id = ?", appID).Updates(updates).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// PUT /api/analysebooks/:id/books/:book_id - изменение м-м связи
func (r *Repository) UpdateBookToApplication(appID, bookID uint, updateData ds.BookToAppl) error {
	var link ds.BookToAppl
	if err := r.db.Where("appl_id = ? AND book_id = ?", appID, bookID).First(&link).Error; err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if updateData.Description != nil {
		updates["description"] = *updateData.Description
	}

	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&link).Updates(updates).Error
}

// func (r *Repository) GetDraftAppl(userID uint) (*ds.Application, error) {
// 	var appl ds.Application

// 	err := r.db.Where("creator_id = ? AND status = ?", userID, ds.StatusDraft).First(&appl).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &appl, nil
// }

// func (r *Repository) CreateAppl(appl *ds.Application) error {
// 	return r.db.Create(appl).Error
// }

// func (r *Repository) AddBookToAppl(applID, bookID uint) error {
// 	var count int64

// 	r.db.Model(&ds.BookToAppl{}).Where("appl_id = ? AND book_id = ?", applID, bookID).Count(&count)
// 	if count > 0 {
// 		return errors.New("book already in application")
// 	}

// 	link := ds.BookToAppl{
// 		ApplID: applID,
// 		BookID: bookID,
// 	}
// 	return r.db.Create(&link).Error
// }

// func (r *Repository) GetApplWithBooks(applID uint) (*ds.Application, error) {
// 	var appl ds.Application

// 	err := r.db.Preload("BooksLink.Book").First(&appl, applID).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	if appl.Status == ds.StatusDeleted {
// 		return nil, errors.New("book page not found or has been deleted")
// 	}

// 	return &appl, nil
// }

// func (r *Repository) LogicallyDeleteAppl(applID uint) error {
// 	result := r.db.Exec("UPDATE applications SET status = ? WHERE id = ?", ds.StatusDeleted, applID)
// 	return result.Error
// }
