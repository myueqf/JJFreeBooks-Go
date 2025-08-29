package main

import (
	"JJFreeBooks/api"
	"JJFreeBooks/config"
	"fmt"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	novelId := 7235273
	chapterId := 55
	chapterDetail, err := api.GetVIPChapterContent(appConfig.Token, novelId, chapterId)
	if err != nil {
		panic(err)
	}
	fmt.Println(chapterDetail)
}
