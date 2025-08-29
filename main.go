package main

import (
	"JJFreeBooks/api"
	"JJFreeBooks/config"
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	version = "dev"     // 默认开发版
	commit  = "none"    // Git 提交哈希
	date    = "unknown" // 构建时间
)

func main() {
	fmt.Println("=======晋江免费小说下载器=======")
	fmt.Println("项目开源地址: https://github.com/MEMLTS/JJFreeBooks-Go")
	fmt.Println("项目作者: MapleLeaf")
	fmt.Println("版本:", version)
	fmt.Println("构建信息:", commit, "@", date)
	fmt.Println("=============================")
	fmt.Println("正在加载配置...")
	appConfig, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println("========加载配置成功========")
	fmt.Printf("Token:%s\n", appConfig.Token)
	fmt.Printf("Cron:%s\n", appConfig.Cron)
	fmt.Println("===========================")
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
			time.Sleep(time.Millisecond * 500)
			// 休眠 500ms
		}
		err = os.WriteFile(bookDir, []byte(content), 0644)
		time.Sleep(time.Second * 2)
		// 休眠 2s
	}
	return true, nil
}
