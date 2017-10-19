package main

import (
	"fmt"
)

type MyType struct {
	Inner *Inner
	Value int
}

type Inner struct {
	Val int
}

func TestPointer(inner Inner) {
	fmt.Println(inner)
}

func Append(array []int, elms ...int) {
	na := append(array, elms...)
	fmt.Printf("value of na %+v \r\n", na)
	array = na
}

func Moidif(array []int, index int, val int) {
	array[index] = val
}

type ArrayContainer struct {
	Elems *[]int
}

func (ac *ArrayContainer) Print() {
	fmt.Printf("%+v \r\n", *ac.Elems)
}

func main() {
	/*var inner = &Inner{}
	inner.Val = 10
	fmt.Println("inner: ", inner)

	inner.Val = 5
	fmt.Println("inner: ", inner)

	mt := MyType{Inner: inner, Value: 1}

	fmt.Println("mt.Inner.Val: ", mt.Inner.Val)

	mt.Inner.Val = 20

	fmt.Println("mt.Inner.Val: ", mt.Inner.Val)

	fmt.Println("inner: ", inner)

	inner.Val = 15
	fmt.Println("inner: ", inner)
	fmt.Println("mt.Inner.Val: ", mt.Inner.Val)

	TestPointer(*inner)*/

	c := make(chan int, 5)
	fmt.Println("c.cap:", cap(c))
	fmt.Println("c.len:", len(c))

	c <- 1
	c <- 2
	fmt.Println("c.cap:", cap(c))
	fmt.Println("c.len:", len(c))

}
