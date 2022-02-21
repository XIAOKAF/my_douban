package main

import (
	"gin/api"
	"gin/dao"
)

func main() {
	dao.InitDB()
	api.InitEngine()
}
