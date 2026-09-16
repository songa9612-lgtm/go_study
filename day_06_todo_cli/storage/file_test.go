package storage

import (
	"go_study/day_06_todo_cli/model"
	"os"
	"testing"
)

func TestFileStorage(t *testing.T) {
	testFile := "test_todo.json"
	s := NewFileStorage(testFile)
	defer os.Remove(testFile)

	fakeItems := []model.Item{
		{
			ID:    1,
			Title: "买牛奶",
			Done:  false,
		},
		{
			ID:    2,
			Title: "写代码",
			Done:  true,
		},
	}

	err := s.Save(fakeItems)
	if err != nil {
		t.Fatalf("保存数据失败: %v", err)
	}

	items, err := s.Load()
	if err != nil {
		t.Fatalf("读取数据失败: %v", err)
	}

	if len(items) != 2 {
		t.Errorf("数目异常")
	}

	// 1. 独立检查第一项的标题
	if items[0].Title != "买牛奶" {
		t.Errorf("第一项标题不一致: 期望 %s, 实际得到 %s", "买牛奶", items[0].Title)
	}

	// 2. 独立检查第二项的完成状态
	if items[1].Done != true {
		t.Errorf("第二项状态不一致: 期望 %v, 实际得到 %v", true, items[1].Done)
	}
}
