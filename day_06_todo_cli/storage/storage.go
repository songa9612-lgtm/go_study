package storage

import "go_study/day_06_todo_cli/model"

type Storage interface {
	// Load 从存储介质中加载所有待办事项
	Load() ([]model.Item, error)

	// Save 将待办切片整体持久化到存储介质中
	Save(items []model.Item) error
}
