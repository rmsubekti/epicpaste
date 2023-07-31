package model

import (
	"epicpaste/system/utils"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type Paste struct {
	ID         string    `gorm:"primarykey:true" json:"id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Public     bool      `json:"public"`
	Languange  string    `json:"language"`
	Tags       *[]Tag    `gorm:"Many2Many:master.paste_tag;FOREIGNKEY:ID;ASSOCIATION_FOREIGNKEY:ID;" json:"tag,omitempty"`
	Category   *Category `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	CategoryId *uint     `json:"-"`
	Paster     User      `gorm:"foreignKey:CreatedBy" json:"creator"`
	CreatedBy  string    `json:"-"`
	CreatedAt  time.Time `time_format:"sql_date" json:"created_at"`
	UpdatedAt  time.Time `time_format:"sql_date" json:"updated_at"`
}

type Pastes []Paste

func (Paste) TableName() string {
	return "master.paste"
}

func (p *Paste) Create() error {
	if len(p.Content) < 1 {
		return errors.New("content must containt at least a word")
	}
	p.ID = uuid.NewString()
	return db.Create(&p).Error
}

func (p *Paste) Update() (err error) {
	var paste Paste
	if len(p.Content) < 1 {
		return errors.New("content must containt at least a word")
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
	return db.Preload(clause.Associations).First(&p, "id = ? ", id).Error
}

func (ps *Pastes) List(paginator *utils.Paginator) (err error) {
	if err = paginator.SetCount(db.Model(&Paste{}).Where("public = ?", true)); err != nil {
		return
	}
	if err = db.Scopes(paginator.Scopes()).Preload(clause.Associations).Find(&ps, "public = ?", true).Error; err != nil {
		return
	}
	paginator.Paginate(ps)
	return
}

func (ps *Pastes) ListByUser(userId string, public bool, paginator *utils.Paginator) (err error) {
	count := db.Model(&Paste{}).Where("created_by = ?", userId)
	pastes := db.Scopes(paginator.Scopes()).Preload(clause.Associations).Where("created_by = ?", userId)

	if public {
		if err = paginator.SetCount(count.Where("public = ?", public)); err != nil {
			return
		}
		if err = pastes.Where("public = ?", public).Find(&ps).Error; err != nil {
			return
		}
	} else {
		if err = paginator.SetCount(count); err != nil {
			return
		}
		if err = pastes.Find(&ps).Error; err != nil {
			return
		}
	}

	paginator.Paginate(pastes)
	return
}
