package ds

import "time"

type AnalyseBooks struct {
	ID               uint       `gorm:"primaryKey;column:id"`
	Status           int        `gorm:"column:status;not null"`
	CreationDate     time.Time  `gorm:"column:creation_date;not null"`
	CreatorID        uint       `gorm:"column:creator_id;not null"` // Внешний ключ
	ModeratorID      *uint      `gorm:"column:moderator_id"`
	FormingDate      *time.Time `gorm:"column:forming_date"`
	ComplitionDate   *time.Time `gorm:"column:complition_date"`
	AwgWordLen       float64    `gorm:"column:awgWordLen;type:numeric(10, 2)"`
	LexicalDiversity float64    `gorm:"column:lexicalDiversity;type:numeric(10, 2)"`
	ConjunctionFreq  float64    `gorm:"column:conjunctionFreq;type:numeric(10, 3)"`
	AvgSentenceLen   float64    `gorm:"column:avgSentenceLen;type:numeric(10, 2)"`

	// --- СВЯЗИ ---
	// Отношение "принадлежит к": каждая заявка принадлежит одному пользователю.
	Creator   Users  `gorm:"foreignKey:CreatorID"`
	Moderator *Users `gorm:"foreignKey:ModeratorID"`
	// Отношение "один-ко-многим" к связующей таблице:
	// У одной заявки может быть много книг.
	BooksLink []BookToAppl `gorm:"foreignKey:ApplID"`
}
