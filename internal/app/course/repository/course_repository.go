package repository

import (
	"context"

	"github.com/valentinusdelvin/studytrack-be/domain/entity"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func New(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) List(ctx context.Context, userID uint64) ([]entity.Course, error) {
	var rows []entity.Course
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("id ASC").Find(&rows).Error
	return rows, err
}
func (r *Repository) Get(ctx context.Context, userID, id uint64) (*entity.Course, error) {
	var row entity.Course
	err := r.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}
func (r *Repository) Create(ctx context.Context, row *entity.Course) error {
	return r.db.WithContext(ctx).Create(row).Error
}
func (r *Repository) Save(ctx context.Context, row *entity.Course) error {
	return r.db.WithContext(ctx).Save(row).Error
}
func (r *Repository) Delete(ctx context.Context, userID, id uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND mata_kuliah_id = ?", userID, id).Delete(&entity.Task{}).Error; err != nil {
			return err
		}
		return tx.Where("user_id = ? AND id = ?", userID, id).Delete(&entity.Course{}).Error
	})
}
