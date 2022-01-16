package models

import (
	"github.com/jinzhu/gorm"
)

type PendingTask struct {
    gorm.Model
    Email  string `gorm:"type:varchar(100);unique_index;not null" json:"email"`
    Task   Task        `gorm:"ForeignKey:TaskID"                  json: "-`
	TaskID uint         `gorm:"not null"                  json:"task"`
}

func (v *PendingTask) Save(db *gorm.DB) (*PendingTask, error) {
    var err error = db.Debug().Create(&v).Error
    if err != nil {
        return &PendingTask{}, err
    }
    return v, nil
}

func GetPendingTaskByEmail(email string, db *gorm.DB) (*PendingTask, error) {
     PendingTask:= &PendingTask{}
    if err := db.Debug().Find(PendingTask, "email = ?", email).Error; err != nil {
        return nil, err
    }
    return PendingTask, nil
}

func DeletePendingTask(id int, db *gorm.DB) error {
    PendingTask:= &PendingTask{}
    if err := db.Debug().Delete(PendingTask, "id = ?", id).Error; err != nil {
        return err
    }
    return nil
}

