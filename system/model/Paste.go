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
	Title      string    `gorm:"type:varchar(125);not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Public     *bool     `json:"public"`
	Languange  string    `json:"language"`
	Tags       *Tags     `gorm:"Many2Many:master.paste_tag;FOREIGNKEY:ID;ASSOCIATION_FOREIGNKEY:ID;" json:"tag,omitempty"`
	Category   *Category `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	CategoryId *uint     `json:"-"`
	Paster     User      `gorm:"foreignKey:CreatedBy" json:"paster" swaggerignore:"true"`
	CreatedBy  string    `json:"-" swaggerignore:"true"`
	CreatedAt  time.Time `time_format:"sql_date" json:"created_at" swaggerignore:"true"`
	UpdatedAt  time.Time `time_format:"sql_date" json:"updated_at" swaggerignore:"true"`
}

type Pastes []Paste

func (Paste) TableName() string {
	return "master.paste"
}

func (p *Paste) Create() error {
	if len(p.Title) < 1 {
		return errors.New("title must containt at least a word")
	}
	if len(p.Content) < 1 {
		return errors.New("content must containt at least a word")
	}

	p.ID = uuid.NewString()
	p.Category.SetAvailable()
	p.Tags.SetAvailable()
	return db.Create(&p).Error
}

func (p *Paste) Update() (err error) {
	var paste Paste
	tagTemp := &Tags{}
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
	p.Category.SetAvailable()
	tagTemp = p.Tags
	if err = db.Model(&p).Association("Tags").Clear(); err != nil {
		return
	}
	tagTemp.SetAvailable()
	p.Tags = tagTemp
	return db.Omit("Paster").Save(&p).Error
}

func (p *Paste) Delete() (err error) {
	var paste Paste

	if err = db.First(&paste, "id = ?", p.ID).Error; err != nil {
		return
	}

	if paste.CreatedBy != p.CreatedBy {
		return errors.New("you dont have access to this paste")
	}

	if err = db.Model(&p).Association("Tags").Clear(); err != nil {
		return
	}
	err = db.Delete(&p).Error
	p = &paste
	return
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

func (ps *Pastes) ListByUser(username string, public bool, paginator *utils.Paginator) (err error) {
	count := db.Model(&Paste{}).Where("created_by = ?", username)
	pastes := db.Scopes(paginator.Scopes()).Preload(clause.Associations).Where("created_by = ?", username)

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

	paginator.Paginate(ps)
	return
}

func (ps *Pastes) ListByCategory(category string, paginator *utils.Paginator) (err error) {
	count := db.Model(&Paste{}).InnerJoins("Category", db.Where(&Category{Name: category})).Where("public = ?", true)
	pastes := db.Scopes(paginator.Scopes()).InnerJoins("Category", db.Where(&Category{Name: category})).Preload(clause.Associations).Where("public = ?", true)
	if err = paginator.SetCount(count); err != nil {
		return
	}
	if err = pastes.Find(&ps).Error; err != nil {
		return
	}
	paginator.Paginate(ps)
	return
}

func (ps *Pastes) ListByTag(tag string, paginator *utils.Paginator) (err error) {
	count := db.Model(&Paste{}).
		Joins("JOIN (?) AS matched ON paste_id = paste.id",
			db.Select("paste_id").
				Table(`"master"."paste_tag" qt`).
				Joins(`JOIN "master"."tag" t ON qt.tag_id = t.id`).
				Where("t.name = ?", tag).
				Group("paste_id"),
		).Where("public = ?", true)
	pastes := db.Scopes(paginator.Scopes()).
		Joins("JOIN (?) AS matched ON paste_id = paste.id",
			db.Select("paste_id").
				Table(`"master"."paste_tag" qt`).
				Joins(`JOIN "master"."tag" t ON qt.tag_id = t.id`).
				Where("t.name = ?", tag).
				Group("paste_id"),
		).Preload(clause.Associations).Where("public = ?", true)

	if err = paginator.SetCount(count); err != nil {
		return
	}
	if err = pastes.Find(&ps).Error; err != nil {
		return
	}
	paginator.Paginate(ps)
	return
}
