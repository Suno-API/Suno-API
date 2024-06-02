package po

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type BaseQuery struct {
	ID      int64
	IDs     []int64
	OrderBy string
	Num     int
	Offset  int
}

func (bq BaseQuery) buildBaseQuery(db *gorm.DB) *gorm.DB {
	if bq.ID != 0 {
		db = db.Where("id = ?", bq.ID)
	}
	if len(bq.IDs) > 0 {
		db = db.Where("id in (?)", bq.IDs)
	}

	if bq.OrderBy != "" {
		db = db.Order(bq.OrderBy)
	}
	if bq.Num != 0 {
		db = db.Limit(bq.Num)
	}
	if bq.Offset != 0 {
		db = db.Offset(bq.Offset)
	}
	return db
}

func RecordExist(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return false, err
}
