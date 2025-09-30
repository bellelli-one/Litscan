package ds

type Books struct {
	ID               uint    `gorm:"primaryKey;column:id"`
	Title            string  `gorm:"column:title;size:255;not null"`
	Text             string  `gorm:"column:text;not null"`
	Image            *string `gorm:"column:image;size:255"`
	AvgWordLen       float64 `gorm:"column:avgWordLen;type:numeric(10, 2)"`
	LexicalDiversity float64 `gorm:"column:lexicalDiversity;type:numeric(10, 2)"`
	ConjunctionFreq  float64 `gorm:"column:conjunctionFreq;type:numeric(10, 3)"`
	AvgSentenceLen   float64 `gorm:"column:avgSentenceLen;type:numeric(10, 2)"`
	Status           *bool   `gorm:"column:status"`

	// --- СВЯЗИ ---
	// Отношение "один-ко-многим" к связующей таблице:
	// Один фактор может быть использован во многих заявках.
	ApplLinks []BookToAppl `gorm:"foreignKey:BookID"`
}
