package main

import (
	"JJFreeBooks/api"
	"fmt"
)

func main() {
	a := api.GetBooksList()
	fmt.Println(a)
}
