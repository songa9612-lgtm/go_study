# Todo-CLI（命令行待办工具）第二部分：业务服务解耦与交互路由规格说明书 (PRD)

> **定位**：Milestone 0 终章冲刺 —— 业务服务解耦与命令行端到端闭环  
> **设计哲学**：三层架构彻底解耦。`main.go` 专职参数校验与路由派发；`service.go` 面向接口契约沉淀核心业务规则；底层复用 `storage.Storage` 持久化引擎。  
> **原则**：纯需求与契约规格驱动，严禁直接提供完整实现代码，由学员独立闭卷编写。

---

## 一、 整体架构与分层设计（为什么拆出 `service.go`？）

按照你的高内聚低耦合设计思想，我们将系统划分为清晰的三层结构：

```text
终端命令输入 (Terminal)
       │ 例如: ./todo add "买牛奶"
       ▼
【接入与路由层: main.go】
       │ 职责：检查 os.Args 长度、转换参数类型 (strconv.Atoi)、switch 派发、终端格式化打印
       ▼ 调用干净业务方法，入参为纯数据（title, id），绝不把 os.Args 传进业务层
【业务契约与实现层: service.go】
       │ 契约：TodoService 接口
       │ 实现：Service 结构体（内部持有 storage.Storage）
       │ 职责：算最大 ID、追加 Item、就地修改 Done 状态、切片删除算法、调用存储写盘
       ▼
【存储持久化层: storage/】
       │ 职责：JSON 文件序列化与磁盘读写（已由第一部分完成并通过单测）
```

### 📁 目录组织结构
```text
day_06_todo_cli/
├── model/
│   └── todo.go          # [数据实体] Item 结构体定义（已完成）
├── storage/
│   ├── storage.go       # [存储契约] Storage 接口（已完成）
│   ├── file.go          # [存储实现] JSON 文件读写引擎（已完成）
│   └── file_test.go     # [存储单测] 自动化回归测试（已通过）
├── service.go           # 👈【本次核心新增】专门实现 add / list / done / del 业务方法与接口契约
├── main.go              # 👈【本次重构】仅保留 30 行纯调度器与命令路由
├── part2_prd.md         # 本任务说明书
└── todos.json           # 运行后持久化存储的待办数据文件
```

---

## 二、 业务服务层契约与实现规范 (`service.go`)

在 `day_06_todo_cli/` 目录下新建 `service.go`，包名同样声明为 `package main`（同一个目录，天然互通，无需跨包导入）。

### 1. 业务接口契约定义 (`TodoService`)
在 `service.go` 顶部定义业务接口，明确待办系统支持的 4 大核心能力契约：
```go
type TodoService interface {
    Add(title string) (model.Item, error)
    List() ([]model.Item, error)
    Done(id int) error
    Delete(id int) error
}
```
> 💡 **面向接口价值**：系统业务能力一目了然。后续无论是重构方法、替换底层存储，还是未来升级为 Gin Web API，直接对照此接口履约。

### 2. 业务结构体与依赖注入
- 定义结构体：
  ```go
  type Service struct {
      storage storage.Storage
  }
  ```
- 编写构造函数：
  ```go
  func NewService(s storage.Storage) *Service
  ```
  *(传入第一步写好的 `storage.FileStorage` 实例，完成存储引擎装配)*

---

### 3. 具体方法实现规格

#### ① `Add(title string) (model.Item, error)`
- **入参**：待办标题 `title`（纯净字符串，不包含任何命令行逻辑）。
- **逻辑步骤**：
  1. 调用 `s.storage.Load()` 获取当前所有待办切片；
  2. **最大 ID 算法**：
     - 初始化 `maxID := 0`；
     - 遍历切片，若 `item.ID > maxID`，则更新 `maxID = item.ID`；
     - 新待办 ID 为 `newID = maxID + 1`（防 ID 碰撞，避免删除后 ID 重复）。
  3. 构造新对象：
     - `ID`: `newID`
     - `Title`: `title`
     - `Done`: `false`
     - `CreatedAt`: `time.Now()`
  4. 追加切片：`items = append(items, newItem)`；
  5. 持久化落盘：调用 `s.storage.Save(items)`；
  6. 返回刚创建的 `newItem` 和 `nil`（若中途出错返回 `model.Item{}, err`）。

#### ② `List() ([]model.Item, error)`
- **逻辑步骤**：
  1. 直接调用 `s.storage.Load()`；
  2. 将获取到的切片和可能发生的 `err` 原样返回给上层展示。

#### ③ `Done(id int) error`
- **入参**：目标待办编号 `id`（上层已转为整数）。
- **逻辑步骤**：
  1. 调用 `s.storage.Load()` 获取所有待办；
  2. 声明布尔标记 `found := false`；
  3. **就地修改（重点避坑）**：
     - 使用下标遍历 `for i := range items`；
     - 当 `items[i].ID == id` 时：
       - `items[i].Done = true`（必须用下标修改，切忌用值拷贝副本修改！）；
       - `found = true`；
       - `break` 提前结束循环。
  4. **未找到防御**：若循环结束 `!found`，使用 `fmt.Errorf("未找到 ID 为 %d 的待办事项", id)` 返回错误；
  5. 找到并修改后，调用 `s.storage.Save(items)` 保存回写磁盘；
  6. 返回 `nil`。

#### ④ `Delete(id int) error`
- **入参**：目标待办编号 `id`。
- **逻辑步骤**：
  1. 调用 `s.storage.Load()` 获取所有待办；
  2. 声明目标下标 `targetIndex := -1`；
  3. 遍历切片查找：当 `items[i].ID == id` 时，记录 `targetIndex = i` 并 `break`；
  4. **未找到防御**：若 `targetIndex == -1`，返回 `fmt.Errorf("未找到 ID 为 %d 的待办事项", id)`；
  5. **切片删除算法（复习 Day 3 切片机制）**：
     - `items = append(items[:targetIndex], items[targetIndex+1:]...)`；
  6. 持久化写回磁盘：`s.storage.Save(items)`；
  7. 返回 `nil`。

---

## 三、 接入与路由层规范 (`main.go`)

在重构后的 `main.go` 中，**不再出现任何算 ID、切片追加、切片删除的底层代码**，它只做一件事：**命令行参数路由与终端展示**。

### 1. 第一道全局门神（防 Panic 崩溃）
```go
if len(os.Args) < 2 {
    printUsage()
    return
}
```

### 2. 依赖装配与初始化
```go
store := storage.NewFileStorage("todos.json")
svc := NewService(store)
```

### 3. 命令分流路由 (`switch os.Args[1]`)

#### 选项 A：`case "add":`
- 检查 `len(os.Args) < 3`，若不足则提示：`错误: 请提供待办内容，例如: todo add "买牛奶"` 并退出；
- 调用 `item, err := svc.Add(os.Args[2])`；
- 若失败打印错误；若成功打印：`[成功] 已添加待办 #%d: %s\n`。

#### 选项 B：`case "list":`
- 调用 `items, err := svc.List()`；
- 若失败打印错误；
- **空切片友好提示**：若 `len(items) == 0`，打印 `暂无待办事项，快用 'todo add' 创建一条吧！` 并返回；
- **格式化表格输出**：
  - 打印表头：`ID    状态    待办内容              创建时间`
  - 打印分隔线：`--------------------------------------------------`
  - 遍历打印每一项：
    - 完成标记：`Done` 为 true 打印 `[x]`，为 false 打印 `[ ]`；
    - 时间格式化：使用固定模板 `item.CreatedAt.Format("2006-01-02 15:04")`；
    - 建议使用制表对齐：`fmt.Printf("%-5d %-7s %-20s %s\n", item.ID, status, item.Title, timeStr)`。

#### 选项 C：`case "done":`
- 检查 `len(os.Args) < 3`；
- 使用 `id, err := strconv.Atoi(os.Args[2])` 安全转换：
  - 若 `err != nil`，提示 `错误: ID 必须为合法整数！`；
- 调用 `err := svc.Done(id)`；
- 若失败打印错误；若成功打印：`[成功] 待办 #%d 已标记为完成！\n`。

#### 选项 D：`case "del":`
- 检查 `len(os.Args) < 3`；
- 使用 `strconv.Atoi(os.Args[2])` 安全转换整型；
- 调用 `err := svc.Delete(id)`；
- 若失败打印错误；若成功打印：`[成功] 待办 #%d 已成功删除！\n`。

#### 选项 E：`default:`
- 提示未知命令，并调用 `printUsage()`。

---

## 四、 关键技术陷阱与排雷指引

1. **同包直接访问原则**：
   `service.go` 和 `main.go` 同在 `day_06_todo_cli` 目录下，第一行都写 `package main`。运行或编译时注意：
   - 终端单文件运行 `go run main.go` 会报找不到 `service.go` 里的东西；
   - **正确命令**：在目录内执行 `go run .`（代表编译运行当前目录下的所有 main 包文件），或者 `go run main.go service.go`！
2. **切片就地修改的副本陷阱**：
   在 `Done` 方法中，切记使用 `items[i].Done = true`，严禁用 `for _, item := range items` 去改 `item.Done`。
3. **切片删除后的持久化**：
   删除元素后务必调 `s.storage.Save(items)`，否则磁盘里的数据没有更新。

---

## 五、 端到端全流程验收清单 (Milestone 0 终验)

在 `day_06_todo_cli` 目录下顺序执行以下全套指令：

```bash
# 1. 验证空列表防崩
go run . list

# 2. 连续添加 3 条任务
go run . add "买牛奶"
go run . add "学习Go接口解耦"
go run . add "练习LeetCode"

# 3. 验证列表格式化制表打印
go run . list

# 4. 标记第 2 条为已完成
go run . done 2

# 5. 再次查看列表，验证状态变为 [x]
go run . list

# 6. 删除第 1 条
go run . del 1

# 7. 再次查看列表，确认第 1 条已除，第 2 和第 3 条顺序完好
go run . list

# 8. 边界输入防御测试（必须优雅报错，严禁发生 panic 崩溃）
go run .
go run . add
go run . done abc
go run . done 999
go run . del 999
go run . unknown
```
全部通过即代表 **Todo-CLI 项目正式交付，Milestone 0 圆满通关！**
