package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Todo struct {
	ID          uint64    `json:"id" gorm:"primary_key"`
	Status      string    `gorm:"size:100;not null;unique" json:"status"`
	Description string    `gorm:"size:255;not null;unique" json:"description"`
	Author      User      `json:"author"`
	AuthorID    uint32    `gorm:"primary_key" json:"author_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Todo) Prepare() {
	p.ID = 0
	p.Status = html.EscapeString(strings.TrimSpace(p.Status))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Todo) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Status == "" {
		err = errors.New("required Status")
		errorMessages["required_status"] = err.Error()
	}
	if p.Description == "" {
		err = errors.New("required Description")
		errorMessages["required_Description"] = err.Error()
	}
	if p.AuthorID < 1 {
		err = errors.New("required Author")
		errorMessages["required_author"] = err.Error()
	}
	return errorMessages
}

func (p *Todo) CreateTodo(db *gorm.DB) (*Todo, error) {
	var err error
	err = db.Debug().Model(&Todo{}).Create(&p).Error
	if err != nil {
		return &Todo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Todo{}, err
		}
	}
	return p, nil
}

func (p *Todo) FindAllTodos(db *gorm.DB) (*[]Todo, error) {
	var err error
	todos := []Todo{}
	err = db.Debug().Model(&Todo{}).Limit(100).Find(&todos).Error
	if err != nil {
		return &[]Todo{}, err
	}
	if len(todos) > 0 {
		for i := range todos {
			err := db.Debug().Model(&User{}).Where("id = ?", todos[i].AuthorID).Take(&todos[i].Author).Error
			if err != nil {
				return &[]Todo{}, err
			}
		}
	}
	return &todos, nil
}

func (p *Todo) FindTodoByID(db *gorm.DB, pid uint64) (*Todo, error) {
	var err error
	err = db.Debug().Model(&Todo{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Todo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Todo{}, err
		}
	}
	return p, nil
}

func (p *Todo) UpdateATodo(db *gorm.DB) (*Todo, error) {

	var err error

	err = db.Debug().Model(&Todo{}).Where("id = ?", p.ID).Updates(Todo{Status: p.Status, Description: p.Description, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Todo{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Todo{}, err
		}
	}
	return p, nil
}

func (p *Todo) DeleteATodo(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Todo{}).Where("id = ? and author_id = ?", p.ID).Take(&Todo{}).Delete(&Todo{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("todo not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
