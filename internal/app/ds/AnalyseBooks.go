package ds

import "time"

type AnalyseBooks struct {
	ID             uint       `gorm:"primaryKey;column:id"`
	Status         int        `gorm:"column:status;not null"`
	CreationDate   time.Time  `gorm:"column:creation_date;not null"`
	CreatorID      uint       `gorm:"column:creator_id;not null"` // Внешний ключ
	ModeratorID    *uint      `gorm:"column:moderator_id"`
	FormingDate    *time.Time `gorm:"column:forming_date"`
	CompletionDate *time.Time `gorm:"column:completion_date"`

	// ИСПРАВЛЕННЫЕ НАЗВАНИЯ КОЛОНОК (snake_case):
	AvgWordLen       float64 `gorm:"column:avg_word_len;type:numeric(10, 2)"`
	LexicalDiversity float64 `gorm:"column:lexical_diversity;type:numeric(10, 2)"`
	ConjunctionFreq  float64 `gorm:"column:conjunction_freq;type:numeric(10, 3)"`
	AvgSentenceLen   float64 `gorm:"column:avg_sentence_len;type:numeric(10, 2)"`

	Response *string `gorm:"column:response;type:text"`

	// --- СВЯЗИ ---
	// Отношение "принадлежит к": каждая заявка принадлежит одному пользователю.
	Creator   Users  `gorm:"foreignKey:CreatorID"`
	Moderator *Users `gorm:"foreignKey:ModeratorID"`
	// Отношение "один-ко-многим" к связующей таблице:
	// У одной заявки может быть много книг.
	BooksLink []BookToAppl `gorm:"foreignKey:ApplID"`
}
