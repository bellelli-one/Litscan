package ds

type BookToAppl struct {
	ID          uint    `gorm:"primaryKey;column:id"`
	ApplID      uint    `gorm:"column:appl_id;not null"`
	BookID      uint    `gorm:"column:book_id;not null"`
	Description *string `gorm:"column:description;type:text"`

	// --- СВЯЗИ ---
	// Отношение "принадлежит к" для каждой из связанных таблиц.
	Appl Application `gorm:"foreignKey:ApplID"`
	Book Books       `gorm:"foreignKey:BookID"`
}
