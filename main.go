package main

import (
	"blog-api/db"
	"blog-api/routes"
)

func main() {
	db.ConnectDB()

  e := routes.Routes()

  e.Logger.Fatal(e.Start(":1323"))
}