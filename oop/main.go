package main

import (
	"fmt"
	"oop/submod"
)

var obj = Object{
	ID:    "1",
	Name:  "First Object",
	Type:  "standard",
	Count: 100,
	Value: 2.5,
}

var objPtr = new(Object)

func main() {
	li1 := []int{0, 1, 2, 3}
	li2 := []int{4, 5, 6, 7}
	li1 = append(li1, li2...)
	fmt.Println(li1)
	fmt.Printf("slice: %v\n", li1)
	fmt.Printf("slice: %+v\n", li1)
	fmt.Printf("slice: %#v\n", li1)
	fmt.Println()

	map1 := map[string]interface{}{"key1": "1", "key2": 2}
	map1["key3"] = 3.0
	fmt.Println(map1)
	fmt.Printf("map: %v\n", map1)
	fmt.Printf("map: %+v\n", map1)
	fmt.Printf("map: %#v\n", map1)
	fmt.Println()

	fmt.Println(obj)
	fmt.Printf("struct: %v\n", obj)
	fmt.Printf("struct: %+v\n", obj)
	fmt.Printf("struct: %#v\n", obj)
	fmt.Println()

	obj.AddOne()
	fmt.Printf("Total Value: %f\n", obj.TotalValue())

	objPtr.SetValue(2.2)
	objPtr.AddOne()
	fmt.Printf("Total Value: %f\n", objPtr.TotalValue())
	fmt.Println()

	str := submod.Quack()
	fmt.Println(str)

	useObject()
}
