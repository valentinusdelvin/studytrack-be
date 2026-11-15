package entity

type User struct {
	ID           uint64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string   `gorm:"column:nama;type:varchar(120);not null" json:"nama"`
	Email        string   `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	PasswordHash string   `gorm:"type:text;not null" json:"-"`
	Courses      []Course `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (User) TableName() string { return "users" }
