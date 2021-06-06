package main

import "fmt"

func app() func(string) string {
	t := "hi"
	c := func(b string) string {
		t = t + " " + b
		return t
	}
	return c
}
func main() {
	a := app()
	b := app()
	a("go")
	fmt.Println(b("all"))
	fmt.Println(a("All"))
}
