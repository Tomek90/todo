package application

import (
	tm "todo/models"

	"github.com/gobuffalo/envy"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	MyDB *gorm.DB
}

//function StartApp starts application by opening connection to database and migrating of models
func (app *Application) StartApp() error {
	databaseURL, err := envy.MustGet("DATABASE_URL")
	if err != nil {
		return err
	}

	app.MyDB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return err
	}

	app.MyDB.AutoMigrate(&tm.ToDo{})

	return nil
}
