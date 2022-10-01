package main

import "fmt"

func takesInterface(object IObject) {
	objStr := fmt.Sprintf("%+v", object)
	fmt.Printf("IObject: %34s\n", objStr)
}

//	func takesIPointer(object *IObject) {
//		fmt.Printf("*IObject:\t%+v\n", object)
//	}
func takesObject(object Object) {
	fmt.Printf("Object:   %+v\n", object)
}
func takesPointer(object *Object) {
	fmt.Printf("*Object: %+v\n", object)
}

func useObject() {
	var obj = Object{}
	var objPtr = new(Object)

	// IObject
	takesInterface(obj)
	takesInterface(&obj)
	takesInterface(objPtr)
	takesInterface(*objPtr)

	// *IObject
	//takesIPointer(*obj)
	//takesIPointer(objPtr)

	// Object
	takesObject(obj)
	takesObject(*objPtr)

	// *Object
	takesPointer(&obj)
	takesPointer(objPtr)
}
