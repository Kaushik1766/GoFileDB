// Package gofiledb provides interface implementation for a filedb
package gofiledb

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"

	"github.com/Kaushik1766/GoFileDB/internal/utils"
)

type Entity interface {
	GetID() string
}

type Repository[T Entity] interface {
	Save(T) error
	GetByParameter(map[string]any) ([]T, error)
	DeleteByParameter(map[string]any) error
}

type FileDBRepository[T Entity] struct {
	mu *sync.Mutex
}

func NewFileDBRepository[T Entity]() *FileDBRepository[T] {
	return &FileDBRepository[T]{mu: &sync.Mutex{}}
}

// Save persist data to db, overwrite if pk already present
func (db *FileDBRepository[T]) Save(data T) error {
	db.mu.Lock()
	defer db.mu.Unlock()
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

// GetByParameter get list of matching values given parameter map
func (db *FileDBRepository[T]) GetByParameter(params map[string]any) ([]T, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	var result []T
	var tmp T

	typ := reflect.TypeOf(tmp).Name()
	fileName := fmt.Sprintf("%v.json", typ)

	if exists := utils.FileExists(fileName); !exists {
		return result, fmt.Errorf("file not found")
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var dbData []T
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		return nil, err
	}

	for _, val := range dbData {
		v := reflect.ValueOf(val)
		for key := range params {
			dbVal := v.FieldByName(key).Interface()
			if dbVal == params[key] {
				result = append(result, val)
			}
		}
	}
	return result, nil
}

// DeleteByParameter delete all matching records given parameter map
func (db *FileDBRepository[T]) DeleteByParameter(params map[string]any) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	var result []T
	var tmp T

	typ := reflect.TypeOf(tmp).Name()
	fileName := fmt.Sprintf("%v.json", typ)

	if exists := utils.FileExists(fileName); !exists {
		return fmt.Errorf("file not found")
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var dbData []T
	err = json.Unmarshal(data, &dbData)
	if err != nil {
		return err
	}

	for _, val := range dbData {
		v := reflect.ValueOf(val)
		for key := range params {
			dbVal := v.FieldByName(key).Interface()
			if dbVal != params[key] {
				result = append(result, val)
			}
		}
	}

	updatedData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	os.WriteFile(fileName, updatedData, 0666)
	return nil
}
