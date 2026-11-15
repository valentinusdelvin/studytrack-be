package database

import (
	"context"
	"errors"
	"time"

	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"github.com/valentinusdelvin/studytrack-be/domain/enums"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user, err := seedUser(tx)
		if err != nil {
			return err
		}
		courses, err := seedCourses(tx, user.ID)
		if err != nil {
			return err
		}
		return seedTasks(tx, user.ID, courses)
	})
}

func seedUser(tx *gorm.DB) (*entity.User, error) {
	const demoEmail = "demo@studytrack.test"
	var user entity.User
	err := tx.Where("lower(email) = ?", demoEmail).First(&user).Error
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte("demo1234"), 10)
	if err != nil {
		return nil, err
	}
	user = entity.User{
		Name:         "Demo Mahasiswa",
		Email:        demoEmail,
		PasswordHash: string(hash),
	}
	if err = tx.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func seedCourses(tx *gorm.DB, userID uint64) (map[uint64]uint64, error) {
	seeds := []entity.Course{
		{
			Name:     "Database Systems",
			Code:     "CS 186",
			Lecturer: "Dr. Andi Wijaya",
		},
		{
			Name:     "Data Structures",
			Code:     "CS 61B",
			Lecturer: "Dr. Sari Lestari",
		},
		{
			Name:     "Software Engineering",
			Code:     "CS 169",
			Lecturer: "Dr. Budi Santoso",
		},
	}
	ids := make(map[uint64]uint64, len(seeds))
	for i, course := range seeds {
		var existing entity.Course
		err := tx.Where("user_id = ? AND name = ?", userID, course.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			course.UserID = userID
			if err = tx.Create(&course).Error; err != nil {
				return nil, err
			}
			existing = course
		} else if err != nil {
			return nil, err
		}
		ids[uint64(i+1)] = existing.ID
	}
	return ids, nil
}

func seedTasks(tx *gorm.DB, userID uint64, courseIDs map[uint64]uint64) error {
	type taskSeed struct {
		course                                uint64
		title, description, deadline, created string
	}
	seeds := []taskSeed{
		{
			1,
			"Lab 4: Query Optimization Report",
			"Laporan hasil eksperimen optimasi query beserta analisis query plan.",
			"2026-07-03T23:59:00",
			"2026-06-15T09:00:00",
		},
		{
			2,
			"Project 3: BSTMap Implementation",
			"Implementasi BSTMap beserta pengujian dasarnya.",
			"2026-07-02T10:00:00",
			"2026-06-16T14:00:00",
		},
		{
			3,
			"Sprint 2: Backend API Integration",
			"Integrasi endpoint backend ke frontend untuk sprint 2.",
			"2026-07-06T23:59:00",
			"2026-06-17T11:00:00",
		},
	}
	for _, item := range seeds {
		var existing entity.Task
		err := tx.Where("user_id = ? AND judul = ?", userID, item.title).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		deadline, err := time.ParseInLocation("2006-01-02T15:04:05", item.deadline, time.Local)
		if err != nil {
			return err
		}
		createdAt, err := time.ParseInLocation("2006-01-02T15:04:05", item.created, time.Local)
		if err != nil {
			return err
		}
		task := entity.Task{
			UserID:      userID,
			CourseID:    courseIDs[item.course],
			Title:       item.title,
			Description: item.description,
			Deadline:    deadline,
			Status:      enums.TaskStatusNotDone,
			CreatedAt:   createdAt,
		}
		if seedErr := tx.Create(&task).Error; seedErr != nil {
			return seedErr
		}
	}
	return nil
}
