package storage

import (
	"encoding/json"
	"errors"
	"go_study/day_06_todo_cli/model"
	"os"
)

type FileStorage struct {
	FilePath string //本地存储路径，例如 "todos.json"
}

func NewFileStorage(filePath string) *FileStorage {
	var path FileStorage
	path.FilePath = filePath
	return &path
}

func (f *FileStorage) Save(items []model.Item) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.FilePath, data, 0644)
}

func (f *FileStorage) Load() ([]model.Item, error) {
	data, err := os.ReadFile(f.FilePath)
	//错误处理
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []model.Item{}, nil
		} else {
			return nil, err
		}
	}

	if len(data) == 0 {
		return []model.Item{}, nil
	}

	var items []model.Item
	err = json.Unmarshal(data, &items)
	return items, err
}
