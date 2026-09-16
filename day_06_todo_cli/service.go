package main

import (
	"fmt"
	"go_study/day_06_todo_cli/model"
	"go_study/day_06_todo_cli/storage"
	"time"
)

type Service struct {
	storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{storage: s}
}

type CliTool interface {
	Add(title string) (model.Item, error)
	List() ([]model.Item, error)
	Done(id int) error
	Delete(id int) error
}

func (s *Service) Add(title string) (model.Item, error) {
	nowSlice, err := s.storage.Load()
	if err != nil {
		return model.Item{}, err
	}

	MaxId := 0
	for _, item := range nowSlice {
		if MaxId < item.ID {
			MaxId = item.ID
		}
	}
	newId := MaxId + 1

	var newItem = model.Item{
		ID:        newId,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}

	nowSlice = append(nowSlice, newItem)

	err = s.storage.Save(nowSlice)
	if err != nil {
		return model.Item{}, err
	}

	return newItem, nil

}

func (s *Service) List() ([]model.Item, error) {
	items, err := s.storage.Load()
	if err != nil {
		return []model.Item{}, err
	}

	return items, nil
}

func (s *Service) Done(id int) error {
	items, err := s.storage.Load()

	found := false

	if err != nil {
		return err
	}

	for i := range items {
		if items[i].ID == id {
			items[i].Done = true
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("未找到 ID 为 %d 的待办事项", id)
	}

	err = s.storage.Save(items)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(id int) error {
	items, err := s.storage.Load()
	if err != nil {
		return err
	}

	targetIndex := -1

	for i := range items {
		if items[i].ID == id {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return fmt.Errorf("未找到 ID 为 %d 的待办事项", id)
	}

	items = append(items[:targetIndex], items[targetIndex+1:]...)
	err = s.storage.Save(items)
	if err != nil {
		return err
	}
	return nil
}
