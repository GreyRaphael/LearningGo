package main

import "fmt"

func main() {
	normalBuilder := getBuilder("normal")
	director := newDirector(normalBuilder)
	normalHouse := director.buildHouse()
	fmt.Printf("%#v\n", normalHouse)

	iglooBuilder := getBuilder("igloo")
	director.setBuilder(iglooBuilder)
	iglooHouse := director.buildHouse()
	fmt.Printf("%#v\n", iglooHouse)
}
