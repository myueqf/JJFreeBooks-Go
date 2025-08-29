package main

import (
	"JJFreeBooks/api"
	"JJFreeBooks/config"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("======加载配置成功======")
	fmt.Printf("Token:%s\n", appConfig.Token)
	fmt.Printf("Cron:%s\n", appConfig.Cron)
	fmt.Println("======================")
	c := cron.New()
	_, err = c.AddFunc(appConfig.Cron, func() {
		_, err := DailyTasks(appConfig)
		if err != nil {
			return
		}
	})
	if err != nil {
		panic(err)
	}

	c.Start()
	defer c.Stop()

	// 阻塞主 goroutine，否则程序会退出
	select {}
}

func DailyTasks(config config.Config) (bool, error) {
	fmt.Println("开始执行定时任务")
	// 获取今日免费小说列表
	bookList, err := api.GetBooksList()
	fmt.Println("今日免费小说列表:", len(bookList.Data.Data))
	if err != nil {
		return false, err
	}
	for _, book := range bookList.Data.Data {
		fmt.Println("开始处理:", book.NovelName)
		dataDir := "data"
		_, err = os.Stat(dataDir)
		if os.IsNotExist(err) {
			err = os.Mkdir(dataDir, 0755)
			if err != nil {
				return false, err
			}
		}
		bookDir := dataDir + "/" + book.NovelName + ".txt"
		_, err = os.Stat(bookDir)
		if os.IsNotExist(err) {
			// 创建空文件
			file, err := os.Create(bookDir)
			if err != nil {
				return false, err
			}
			_ = file.Close()
		}
		// 获取章节列表
		chapterList, err := api.GetChapterList(book.NovelID)
		if err != nil {
			return false, err
		}
		var content string
		// 获取章节内容
		for _, chapter := range chapterList.Chapterlist {
			var chapterContent api.ChapterDetail
			if chapter.IsVip == 0 {
				chapterContent, err = api.GetChapterContent(book.NovelID, chapter.ChapterID)
			} else {
				chapterContent, err = api.GetVIPChapterContent(config.Token, book.NovelID, chapter.ChapterID)
			}
			if err != nil {
				return false, err
			}
			content += "第" + chapterContent.ChapterID + "章" + chapterContent.ChapterName + "\n" + chapterContent.Content + "\n"
		}
		err = os.WriteFile(bookDir, []byte(content), 0644)
	}
	return true, nil
}
