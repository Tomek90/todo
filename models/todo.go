package models

import (
	"time"

	"gorm.io/gorm"
)

type ToDo struct {
	gorm.Model
	DateAndTimeOfExpiry *time.Time
	Title               string
	Description         string
	CompletePercentage  int
}

func (todo *ToDo) Create(db *gorm.DB) error {
	return db.Create(todo).Error
}

func (todo *ToDo) Delete(db *gorm.DB) error {
	return db.Delete(todo).Error
}

func (todo *ToDo) Update(db *gorm.DB) error {
	return db.Save(todo).Error
}

func GetToDoByID(db *gorm.DB, id uint) (*ToDo, error) {
	toDo := &ToDo{}

	err := db.First(toDo, id).Error
	if err != nil {
		return nil, err
	}

	return toDo, nil
}

func GetAllToDos(db *gorm.DB) ([]ToDo, error) {
	toDos := []ToDo{}

	err := db.Find(&toDos).Error
	if err != nil {
		return nil, err
	}

	return toDos, nil
}
