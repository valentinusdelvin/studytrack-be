package repository

import (
	"context"

	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

type TaskWithCourse struct {
	entity.Task
	CourseName string `gorm:"column:course_name"`
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) GetTasksByUserID(ctx context.Context, userID uint64) ([]TaskWithCourse, error) {
	var rows []TaskWithCourse
	err := r.db.WithContext(ctx).Table("tasks").
		Select("tasks.*, courses.name AS course_name").
		Joins("JOIN courses ON courses.id = tasks.mata_kuliah_id").
		Where("tasks.user_id = ? AND courses.user_id = ?", userID, userID).
		Order("tasks.created_at DESC").Scan(&rows).Error
	return rows, err
}

func (r *TaskRepository) GetTaskByIDAndUserID(ctx context.Context, userID, id uint64) (*entity.Task, error) {
	var row entity.Task
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *TaskRepository) IsCourseOwnedByUser(ctx context.Context, userID, courseID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Course{}).Where("id = ? AND user_id = ?", courseID, userID).Count(&count).Error
	return count > 0, err
}

func (r *TaskRepository) CreateTask(ctx context.Context, row *entity.Task) error {
	return r.db.WithContext(ctx).Create(row).Error
}
func (r *TaskRepository) UpdateTask(ctx context.Context, row *entity.Task) error {
	return r.db.WithContext(ctx).Save(row).Error
}
func (r *TaskRepository) DeleteTask(ctx context.Context, userID, id uint64) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).Delete(&entity.Task{}).Error
}
