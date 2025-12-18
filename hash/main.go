package main

import (
	"fmt"

	"backend/app/utils"
)

func main() {
	hash, err := utils.HashPassword("ciko123")
	if err != nil {
		panic(err)
	}

	fmt.Println(hash)
}
