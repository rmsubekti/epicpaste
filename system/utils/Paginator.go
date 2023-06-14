package utils

import (
	"math"

	"gorm.io/gorm"
)

type Paginator struct {
	Limit      int         `form:"limit, omitempty" json:"limit"`
	Page       int         `form:"page, omitempty" json:"page"`
	Sort       string      `form:"sort, omitempty" json:"sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
	offset     int
}

func (p *Paginator) Paginate(model any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.offset).Limit(p.Limit).Order(p.Sort)
	}
}

func (p *Paginator) CalculateOffset(total int64) {
	if p.Limit <= 0 {
		p.Limit = 10
	}

	if p.Page <= 0 {
		p.Page = 1
	}

	p.offset = (p.Page - 1) * p.Limit

	if p.Sort == "asc" {
		p.Sort = "id asc"
	} else {
		p.Sort = "id desc"
	}

	p.TotalRows = total
	p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.Limit)))

}
