package main

import (
	"sample/internal/router"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router.Run()
}
