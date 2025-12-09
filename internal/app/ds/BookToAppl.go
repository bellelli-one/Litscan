package ds

type BookToAppl struct {
	ID          uint     `gorm:"primaryKey;column:id"`
	ApplID      uint     `gorm:"column:appl_id;not null"`
	BookID      uint     `gorm:"column:book_id;not null"`
	Description *string  `gorm:"column:description;type:text"`
	Similarity  *float64 `gorm:"column:similarity;type:numeric(5, 4)"`

	// --- СВЯЗИ ---
	// Отношение "принадлежит к" для каждой из связанных таблиц.
	Appl AnalyseBooks `gorm:"foreignKey:ApplID"`
	Book Books        `gorm:"foreignKey:BookID"`
}
