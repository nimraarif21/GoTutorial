package models

import (
	"errors"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Task struct {
    gorm.Model
    Name        string `gorm:"size:100;not null;unique"     json:"name"`
    Description string `gorm:"not null"                     json:"description"`
    DueDate     time.Time `gorm:"not null"                  json:"time"`
    CreatedBy   User   `gorm:"foreignKey:UserID;"           json:"-"`
    UserID      uint   `gorm:"not null"                     json:"user_id"`
    AssignedTo   User `gorm:"ForeignKey:AssignedUserID"    json: "-"`
    AssignedUserID uint `gorm:"not null"                    json:"assigned_user_id"`
}

func (v *Task) Prepare() {
    v.Name = strings.TrimSpace(v.Name)
    v.Description = strings.TrimSpace(v.Description)
	v.DueDate = v.DueDate
    v.CreatedBy = User{}
}

func (v *Task) Validate() error {
    if v.Name == "" {
        return errors.New("name is required")
    }
    if v.Description == "" {
        return errors.New("description about Task is required")
    }
		// if v.DueDate.Before(time.Now()){
		// 	return errors.New(time.Now().String())
		// }
    return nil
}

func (v *Task) Save(db *gorm.DB) (*Task, error) {

    // Debug a single operation, show detailed log for this operation
    var err error = db.Debug().Create(&v).Error
    if err != nil {
        return &Task{}, err
    }
    return v, nil
}

func (v *Task) GetTask(db *gorm.DB) (*Task, error) {
    Task := &Task{}
    if err := db.Debug().Table("Tasks").Where("name = (?)", v.Name).First(Task).Error; err != nil {
        return nil, err
    }
		
    return Task, nil
}

func GetTasks(userID uint, db *gorm.DB) (*[]Task, error) {
    Tasks := []Task{}

    if err := db.Debug().Find(&Tasks, "user_id = ?", userID).Error; err != nil {
        return &[]Task{}, err
    }
    return &Tasks, nil
}

func GetTaskById(id int, db *gorm.DB) (*Task, error) {
    Task := &Task{}
    if err := db.Debug().Find(Task, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return Task, nil
}

func (v *Task) UpdateTask(id int, db *gorm.DB) (*Task, error) {
    UpdatedTask := &Task{}
    if err := db.Debug().Find(UpdatedTask, "id = ?", id).Updates(Task{
        Name:        v.Name,
        Description: v.Description,
        DueDate:    v.DueDate,
        AssignedTo:    v.AssignedTo}).Error; err != nil {
        return UpdatedTask, err
    }
    return v, nil
}

func DeleteTask(id int, db *gorm.DB) error {
    Task := &Task{}
    if err := db.Debug().Delete(Task, "id = ?", id).Error; err != nil {
        return err
    }
    return nil
}


func GetAlltasks(userID uint, sortby string, limit int, offset int, db *gorm.DB) (*[]Task, error) {
    Tasks := []Task{}
    if err := db.Debug().Limit(limit).Offset(offset).Order((sortby)).Find(&Tasks, "user_id = ?", userID).Error; err != nil {
        return &[]Task{}, err
    }
    return &Tasks, nil
}





