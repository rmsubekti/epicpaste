package model

import (
	"epicpaste/system/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

const (
	CODE string = "code"
	NOTE string = "note"
	BLOG string = "blog"
)

type Paste struct {
	ID       string    `gorm:"type:varchar(40);primarykey:true;not null;unique" json:"id"`
	Title    string    `gorm:"type:varchar(70);not null" json:"title"`
	Text     string    `gorm:"type:varchar(254);not null" json:"text"`
	Public   bool      `json:"public"`
	Type     string    `json:"type"`
	Tag      *[]Tag    `gorm:"Many2Many:master.paste_tag;FOREIGNKEY:ID;ASSOCIATION_FOREIGNKEY:ID;" json:"tag,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	//allow nil *uint
	CategoryId *uint       `json:"-"`
	CodeDetail *CodeDetail `gorm:"foreignKey:ID" json:"code_detail,omitempty"`
	Creator    User        `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedBy  string      `json:"-"`
	CreatedAt  time.Time   `time_format:"sql_date" json:"created_at"`
	UpdatedAt  time.Time   `time_format:"sql_date" json:"updated_at"`
}

type Pastes []Paste

func (Paste) TableName() string {
	return "master.paste"
}

func (p *Paste) Create() error {
	if len(p.Title) < 1 {
		p.Title = "Untitled"
	}
	if len(p.Text) < 1 {
		return errors.New("text content must containt at least a word")
	}
	if len(p.Type) < 1 {
		p.Type = NOTE
	}
	if p.Type != CODE {
		p.CodeDetail = nil
	}
	p.ID = uuid.NewString()
	return db.Create(&p).Error
}

func (p *Paste) Update() (err error) {
	var paste Paste
	if len(p.Text) < 1 {
		return errors.New("text content must containt at least a word")
	}

	if err = db.First(&paste, "id = ?", p.ID).Error; err != nil {
		return
	}

	if paste.CreatedBy != p.CreatedBy {
		return errors.New("you dont have access to this paste")
	}

	p.UpdatedAt = time.Now()
	return db.Save(&p).Error
}

func (p *Paste) Delete() (err error) {
	var paste Paste

	if err = db.First(&paste, "id = ?", p.ID).Error; err != nil {
		return
	}

	if paste.CreatedBy != p.CreatedBy {
		return errors.New("you dont have access to this paste")
	}
	p = &paste
	return db.Delete(&p).Error
}

func (p *Paste) Get(id string) error {
	return db.Preload("CodeDetail.App").Preload("CodeDetail.Language").Preload("CodeDetail.OS").Preload(clause.Associations).First(&p, "id = ? ", id).Error
}

func (ps *Pastes) List(paginator *utils.Paginator) (err error) {
	var totalRows int64
	if err = db.Model(&Paste{}).Where("public = ?", true).Count(&totalRows).Error; err != nil {
		return
	}
	paginator.CalculateOffset(totalRows)

	if err = db.Scopes(paginator.Paginate(ps)).Preload(clause.Associations).Find(&ps, "public = ?", true).Error; err != nil {
		return
	}

	paginator.Rows = ps
	return
}

func (ps *Pastes) ListByUser(userId string, public bool, paginator *utils.Paginator) (err error) {
	var totalRows int64
	count := db.Model(&Paste{}).Where("created_by = ?", userId)
	pastes := db.Scopes(paginator.Paginate(ps)).Preload(clause.Associations).Where("created_by = ?", userId)

	if public {
		if err = count.Where("public = ?", public).Count(&totalRows).Error; err != nil {
			return
		}

		paginator.CalculateOffset(totalRows)

		if err = pastes.Where("public = ?", public).Find(&ps).Error; err != nil {
			return
		}
	} else {
		if err = count.Count(&totalRows).Error; err != nil {
			return
		}

		paginator.CalculateOffset(totalRows)

		if err = pastes.Find(&ps).Error; err != nil {
			return
		}
	}

	paginator.Rows = ps
	return
}
