package main

import (
	"errors"
	"fmt"
)

// Go 语言要求：既然你声明了要返回 2 个东西（一个数字，一个 error），
// 那么无论哪条 return 分支，都必须实打实地给出这两个东西！
func safeDivide(a, b int) (float64, error) {

	if b == 0 {
		return 0, errors.New("除数不能为0")
	}

	return float64(a) / float64(b), nil
}

// getStats 接收任意数量的整数，计算并返回最大值、最小值、平均值
func getStats(numbers ...int) (max int, min int, avg float64) {
	// 1. 防御：如果什么参数都没传，直接返回 0
	if len(numbers) == 0 {
		return 0, 0, 0
	}

	// 2. 初始化：以第一个元素作为初始最值
	max = numbers[0]
	min = numbers[0]
	sum := 0

	// 3. 遍历所有传入的数字
	for _, num := range numbers {
		if num > max {
			max = num
		}
		if num < min {
			min = num
		}
		sum += num // 累加求和
	}

	// 4. 计算平均值（注意：必须显式强转为 float64，否则整数相除会丢弃小数！）
	avg = float64(sum) / float64(len(numbers))

	// 5. 返回结果（因为是命名返回值，也可以直接写 return）
	return max, min, avg
}

func main() {
	result, err := safeDivide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
		return // 直接 return，结束 main 函数的执行
	}
	fmt.Println("Result:", result)

	max, min, avg := getStats(1, 2, 3, 4, 5)
	fmt.Printf("Max: %v, Min: %v, Avg: %v\n", max, min, avg)

}
