package pagination

import (
	"context"
	"gorm.io/gorm"
)

func Paginate(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		limit := ctx.Value(Limit).(int64)
		page := ctx.Value(Page).(int64)

		offset := (page - 1) * limit

		return db.Offset(int(offset)).Limit(int(limit))
	}
}
