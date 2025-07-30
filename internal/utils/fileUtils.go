package utils

import (
	"encoding/json"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func SaveData[T any](name string, data T, perm os.FileMode) error {
	marshalData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	err = os.WriteFile(name, marshalData, perm)
	return err
}
