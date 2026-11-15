package entity

import (
	"time"

	"github.com/valentinusdelvin/studytrack-be/domain/enums"
)

type TaskWithCourse struct {
	ID          uint64           `gorm:"primaryKey;autoIncrement"`
	UserID      uint64           `gorm:"not null;index"`
	CourseID    uint64           `gorm:"column:mata_kuliah_id;not null;index"`
	Title       string           `gorm:"column:judul;type:varchar(255);not null"`
	Description string           `gorm:"column:deskripsi;type:text;not null;default:''"`
	Deadline    time.Time        `gorm:"type:timestamp without time zone;not null"`
	Status      enums.TaskStatus `gorm:"type:varchar(16);not null;default:'belum'"`
	CreatedAt   time.Time        `gorm:"not null;autoCreateTime"`
	CourseName  string           `gorm:"column:course_name"`
}
