package main

import "fmt"

func main() {
	ak47, _ := getGun("ak47")
	fmt.Printf("Gun: %#v\n", ak47)
	printDetails(ak47)

	musket, _ := getGun("musket")
	fmt.Printf("Gun: %#v\n", musket)
	printDetails(musket)
}

func printDetails(g IGun) {
	fmt.Printf("Gun: %s, Power: %d\n", g.getName(), g.getPower())
}
