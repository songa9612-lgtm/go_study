// 变量定义了必须使用,不然没法编译
package main

import "fmt"

const (
	_ = iota //第一个iota被赋值为0，但因为被_接收，所以没有被使用
	Pending
	Paid
	Shipping
	Completed
	Cancelled
)

const (
	Monday = iota
	Tuesday
	Wenesday
	Thusday
	Firday
	Satuday
	Sunday
)

func main() {

	//变量声明
	var ( //当多个变量类型相同时，可以省略类型的声明，必须使用小括号将多个变量声明包裹起来
		age    int     = 10
		score  float64 = 9.5
		name   string  = "pual"
		raally bool    = false
	)
	fmt.Println(age, score, name, raally)

	//零值
	var Zero_Value int
	fmt.Println(Zero_Value)

	//类型推导
	Type_deduction := 10
	fmt.Println(Type_deduction)

	// 显式类型转换
	var x float64 = 3.12
	var y float32 = 4.6
	fmt.Println(x + float64(y))

}
