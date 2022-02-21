package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/douban?charset=utf8mb4&parseTime=True")
	if err != nil {
		fmt.Println("failed", err)
		panic(err)
	}
	DB = db
}
