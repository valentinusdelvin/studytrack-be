package entity

type Course struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   uint64 `gorm:"not null;index" json:"user_id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Code     string `gorm:"column:kode;type:varchar(64);not null;default:''" json:"kode"`
	Lecturer string `gorm:"column:dosen;type:varchar(255);not null;default:''" json:"dosen"`
	Tasks    []Task `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (Course) TableName() string { return "courses" }
