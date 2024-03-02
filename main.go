package main

import (
	_ "github.com/lib/pq"
)

// CONNECT DB AND INIT..
func main() {
	database.setupDB()
}
