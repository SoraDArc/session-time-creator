package main

import (
	"coursach/seed"

	_ "github.com/lib/pq"
)

func main() {
	//seed.Init()
	//api.Start()
	seed.ToExcelFile()
}
