package main

import (
	"fmt"
	"github.com/van-pelt/quotes/pkg/quotes"
)

func main() {
	dd := quotes.NewQuotes()
	dd.AddQuotes("dfdf", "dfdfdf", "dfdfdfdfdf")
	fmt.Println(dd.GetAllQuotes())

}
