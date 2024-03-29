package models

import (
	"errors"
	"html"
	"strings"
	"time"

	
	"github.com/jinzhu/gorm"
)

//Todo
type Todo struct {
	ID          uint64    `json:"id" gorm:"primary_key"`
	Status      string    `gorm:"size:100;not null;unique" json:"status"`
	Description string    `gorm:"size:255;not null;unique" json:"description"`
	Author      User      `json:"author"`
	AuthorID    uint32    `gorm:"not null" json:"author_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}



func (p *Todo) Prepare() {
	
	p.Status = html.EscapeString(strings.TrimSpace(p.Status))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}


//Validate
func (p *Todo) Validate() map[string]string {

	var err error

	var errorMessages = make(map[string]string)

	if p.Status == "" {
		err = errors.New("required Status")
		errorMessages["Required_Status"] = err.Error()
	}
	if p.Description == "" {
		err = errors.New("required Description")
		errorMessages["Required_Description"] = err.Error()
	}
	if p.AuthorID < 1 {
		err = errors.New("required Author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}
//CreateTodo
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
	return p , nil
}

//FindAllTodos
func (p *Todo) FindTodosByUserID(db *gorm.DB, userID uint32) (*[]Todo, error) {
    var err error
    todos := []Todo{}
    err = db.Debug().Model(&Todo{}).Where("author_id = ?", userID).Limit(100).Find(&todos).Error
    if err != nil {
        return &[]Todo{}, err
    }
    return &todos, nil
}
//FindTodoByID
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
//UpdateATodo
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
	return &Todo{}, nil
}
//DeleteATodo
func (p *Todo) DeleteATodo(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Todo{}).Where("id = ? and author_id = ?", pid, uid).Take(&Todo{}).Delete(&Todo{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("todo not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (c *Todo) DeleteUserTodos(db *gorm.DB, uid uint32) (int64, error) {
	todos := []Todo{}
	db = db.Debug().Model(&Todo{}).Where("author_id = ?", uid).Find(&todos).Delete(&todos)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}