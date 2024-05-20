package fact

import (
	"context"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db.Table("facts"),
	}
}

// Insert Сохраняем в БД факты
func (r *Repository) Insert(ctx context.Context, col Collection) error {
	if len(col) < 1 {
		return nil
	}

	return r.db.WithContext(ctx).Create(col).Error
}

// GetFirstUnsentFact для получения фактов без indicator_to_mo_fact_id
func (r *Repository) GetFirstUnsentFact(ctx context.Context, fact *Fact) error {
	return r.db.WithContext(ctx).
		Where("indicator_to_mo_fact_id = ?", 0).
		Limit(1).
		Find(fact).Error
}

// UpdateFact для обновления факта в БД
func (r *Repository) UpdateFact(ctx context.Context, fact *Fact) error {
	return r.db.WithContext(ctx).
		Save(fact).
		Error
}
