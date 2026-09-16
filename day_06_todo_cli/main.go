package main

import (
	"fmt"
	"go_study/day_06_todo_cli/storage"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	store := storage.NewFileStorage("todo.json")
	svc := NewService(store)

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("错误: 请提供待办内容，例如: todo add '买牛奶'")
			printUsage()
			return
		}
		items, err := svc.Add(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[成功] 已添加待办 #%d: %s\n", items.ID, items.Title)

	case "list":
		items, err := svc.List()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(items) == 0 {
			fmt.Println("暂无待办事项，快用 'todo add' 创建一条吧!")
			return
		}

		fmt.Println("ID 状态 待办内容 创建时间")
		fmt.Println("--------------------------------------------------")
		for _, item := range items {
			timeStr := item.CreatedAt.Format("2006-01-02 15:04:05")
			if item.Done == false {
				fmt.Printf("%-5d %-7s %-20s %s\n", item.ID, "[ ]", item.Title, timeStr)
			} else {
				fmt.Printf("%-5d %-7s %-20s %s\n", item.ID, "[x]", item.Title, timeStr)
			}
		}

	case "done":
		if len(os.Args) < 3 {
			fmt.Println("错误,请提供具体命令")
			printUsage()
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("id必须为合法整数")
			return
		}

		err = svc.Done(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("[成功] 待办 #%d 已标记为完成！\n", id)

	case "del":
		if len(os.Args) < 3 {
			fmt.Println("错误,请提供具体命令")
			printUsage()
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("id必须为合法整数")
			return
		}

		err = svc.Delete(id)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("[成功] 待办 #%d 已成功删除！\n", id)

	default:
		fmt.Println("未知命令")
		printUsage()
	}

}

func printUsage() {
	fmt.Println("=== 欢迎使用 Todo-CLI 待办管理工具 ===")
	fmt.Println("使用方式 (Usage):")
	fmt.Println("  todo add <任务内容>    - 添加一条新的待办事项")
	fmt.Println("  todo list             - 查看所有待办事项列表")
	fmt.Println("  todo done <任务ID>     - 标记指定 ID 的待办为已完成")
	fmt.Println("  todo del <任务ID>      - 删除指定 ID 的待办事项")
	fmt.Println("示例:")
	fmt.Println("  go run . add \"买牛奶\"")
	fmt.Println("  go run . done 1")
}
