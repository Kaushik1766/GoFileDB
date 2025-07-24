// Package gofiledb provides interface implementation for a filedb
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/Kaushik1766/GoFileDB/utils"
)

type Entity interface {
	GetID() string
}

type Repository[T Entity] interface {
	// Save persist data to db, overwrite if pk already present
	Save(T) error
	// GetByParameter get top value given the required parameter map
	GetByUniqueParameter(map[string]any) (T, error)
	// GetByParameter get list of matching values given parameter map
	GetByParameter(map[string]any) ([]T, error)
	// DeleteByParameter delete all matching records given parameter map
	DeleteByParameter(map[string]any) error
}

type FileDBRepository[T Entity] struct{}

func (db *FileDBRepository[T]) Save(data T) error {
	typ := reflect.TypeOf(data).Name()
	fileName := fmt.Sprintf("%v.json", typ)

	if exists := utils.FileExists(fileName); !exists {
		var tmp []T
		err := utils.SaveData(fileName, tmp, 0666)
		if err != nil {
			return err
		}
	}

	prevDataBytes, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	var prevData []T

	if err := json.Unmarshal(prevDataBytes, &prevData); err != nil {
		return err
	}

	updated := false

	for i, val := range prevData {
		if val.GetID() == data.GetID() {
			prevData[i] = data
			updated = true
			break
		}
	}

	if !updated {
		prevData = append(prevData, data)
	}

	err = utils.SaveData(fileName, prevData, 0666)

	return err
}

type User struct {
	Id   string
	Name string
	Age  int
}

func (u User) GetID() string {
	return u.Id
}

func main() {
	db := FileDBRepository[User]{}
	db.Save(User{"adsfaddsdf", "kaushik", 41})
}
