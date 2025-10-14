package ds

type Users struct {
	ID        uint   `gorm:"primaryKey;column:id"`
	Username  string `gorm:"unique;column:username;size:255;not null"`
	Password  string `gorm:"unique;column:password;size:255;not null"`
	Moderator bool   `gorm:"column:moderator;not null"`
	// --- СВЯЗИ ---
	// Отношение "один-ко-многим": один пользователь может иметь много поисковых сессий.
	AnalyseBooks []AnalyseBooks `gorm:"foreignKey:CreatorID"`
}
