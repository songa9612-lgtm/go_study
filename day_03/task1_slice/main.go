package main

import "fmt"

func main() {
	s := make([]int, 0, 2) //make创建空切片,长度为0,容量为2
	for i := range 10 {
		s = append(s, i)
		fmt.Println(s, len(s), cap(s), &s[i])
	}

	origin := []int{10, 20, 30, 40, 50}

	//切片的底层数组是共享的,所以修改sub切片的元素会影响origin切片的元素
	sub := origin[0:4] //sub{10,20,30,40}
	sub[0] = 999       //origin{999,20,30,40,50}
	fmt.Println(sub, origin)

	origin[0] = 10
	s2 := make([]int, len(origin), cap(origin))
	//copy函数的作用是将origin切片的内容复制到s2切片中,复制依据的元素个数s2的len而非cap
	copy(s2, origin)
	s2[0] = 888
	fmt.Println(s2, origin)

	origin[0] = 10
	origin = removeAT(origin, 2)
	fmt.Println(origin, len(origin), cap(origin))
}

func removeAT(s []int, index int) []int {

	//...是“解包运算符”（相当于把一个切片“打散”成一个个独立的元素）。
	return append(s[:index], s[index+1:]...)
}
