package seed

import (
	"coursach/models"
	"errors"
	"fmt"
)

// Saves data from files.
func SeedObjects(name string) {
	objects, tableName, err := LoadObjects(name)
	if err != nil {
		fmt.Println(errors.New("seed/seeding.go:SeedObjects: не удалось загрузить данные"))
	}
	for _, obj := range objects {
		_, err := models.Save(tableName, obj)
		if err != nil {
			panic(errors.New("coursach/seed:SeedObjects: Ошибка в сохранении"))
		}
	}
}
