package main

import "fmt"

type aStruct struct {
	a int
	b int
}

func (as *aStruct) PrintA() {
	fmt.Printf("%v\n", as.a)
}

func (as aStruct) CopyPrintA() {
	fmt.Printf("%v\n", as.a)
}


func testFunc1() {
	v := aStruct{a : 1, b:2}
	v.PrintA()
	defer v.PrintA()
	v.a = 100
}

func testCopyFunc1() {
	v := aStruct{a : 1, b:2}
	v.PrintA()
	defer v.CopyPrintA()
	v.a = 100
}

func testFunc2() {
	v := aStruct{a : 1, b:2}
	v.PrintA()
	defer v.PrintA()
	v = aStruct{a : 3, b:4}
}

func testCopyFunc2() {
	v := aStruct{a : 1, b:2}
	v.PrintA()
	defer v.CopyPrintA()
	v = aStruct{a : 3, b:4}
}

func main() {
	testFunc1()
	fmt.Println("-----")
	testCopyFunc1()
	fmt.Println("-----")
	testFunc2()
	fmt.Println("-----")
	testCopyFunc2()
	fmt.Println("-----")
}
