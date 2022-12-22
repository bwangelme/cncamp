package main

import (
	"fmt"
	"reflect"
)

type T struct {
	A int
}

func (t T) Add() int {
	return t.A + 42
}

func main() {
	t := T{20}

	v := reflect.ValueOf(t)
	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("Field %d %v\n", i, v.Field(i))
	}
	for i := 0; i < v.NumMethod(); i++ {
		fmt.Printf("Method %d %v\n", i, v.Method(i))
	}

	res := v.Method(0).Call([]reflect.Value{})
	fmt.Println(res[0])
}
