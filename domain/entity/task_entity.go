package entity

import (
	"time"

	"github.com/valentinusdelvin/studytrack-be/domain/enums"
)

type Task struct {
	ID          uint64           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint64           `gorm:"not null;index" json:"user_id"`
	CourseID    uint64           `gorm:"column:mata_kuliah_id;not null;index" json:"mata_kuliah_id"`
	Title       string           `gorm:"column:judul;type:varchar(255);not null" json:"judul"`
	Description string           `gorm:"column:deskripsi;type:text;not null;default:''" json:"deskripsi"`
	Deadline    time.Time        `gorm:"type:timestamp without time zone;not null" json:"deadline"`
	Status      enums.TaskStatus `gorm:"type:varchar(16);not null;default:'belum'" json:"status"`
	CreatedAt   time.Time        `gorm:"not null;autoCreateTime" json:"created_at"`
	User        User             `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Course      Course           `gorm:"foreignKey:CourseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}

func (Task) TableName() string { return "tasks" }
