package main

import (
	"JJFreeBooks/api"
	"JJFreeBooks/config"
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

// 全局变量声明区 🌍
var (
	version = "dev"     // 默认开发版 🛠️ - 表示当前是开发版本
	commit  = "none"    // Git 提交哈希 🔖 - 源代码版本控制标识
	date    = "unknown" // 构建时间 ⏰ - 程序编译打包的时间
)

// 主函数 - 程序入口点 🚀
func main() {
	// 炫酷的启动横幅 🎉
	fmt.Println("✨=======晋江免费小说下载器=======✨")
	fmt.Println("📖 项目开源地址: https://github.com/MEMLTS/JJFreeBooks-Go")
	fmt.Println("👨‍💻 项目作者: MapleLeaf 🍁")
	fmt.Println("🏷️ 版本:", version)
	fmt.Println("🔧 构建信息:", commit, "@", date)
	fmt.Println("⏰ 启动时间:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("=====================================")

	// 加载配置 🗂️
	fmt.Println("🔄 正在加载配置...")
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println("❌ 配置加载失败:", err)
		panic("🔥 配置文件加载失败，请检查config文件是否存在且格式正确！")
	}

	fmt.Println("✅ ========加载配置成功========")
	fmt.Printf("🔑 Token:%s\n", appConfig.Token)
	fmt.Printf("⏰ Cron表达式:%s\n", appConfig.Cron)
	fmt.Println("===============================")

	// 创建cron调度器 ⏲️
	fmt.Println("🔄 初始化定时任务调度器...")
	c := cron.New()

	// 添加定时任务 📅
	fmt.Printf("🎯 添加定时任务，表达式: %s\n", appConfig.Cron)
	_, err = c.AddFunc(appConfig.Cron, func() {
		fmt.Printf("⏰ 定时任务触发于: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		success, err := DailyTasks(appConfig)
		if err != nil {
			fmt.Printf("❌ 定时任务执行失败: %s\n", err)
		} else if success {
			fmt.Println("✅ 定时任务执行完成!")
		}
	})

	if err != nil {
		fmt.Println("❌ 添加定时任务失败:", err)
		panic("🔥 Cron表达式可能无效，请检查配置！")
	}

	fmt.Println("✅ 定时任务添加成功!")
	fmt.Println("🚀 启动定时任务调度器...")
	c.Start()
	defer c.Stop() // 优雅关闭 🔄

	fmt.Println("🌈 程序已启动并运行中...")
	fmt.Println("💡 提示: 按Ctrl+C可退出程序")
	fmt.Println("=====================================")

	// 阻塞主 goroutine，否则程序会退出 ⛔
	select {} // 无限阻塞，保持程序运行 ♾️
}

// DailyTasks 每日任务处理函数 📋
// 参数: config - 应用程序配置
// 返回值: bool - 任务是否成功, error - 错误信息
func DailyTasks(config config.Config) (bool, error) {
	fmt.Println("——————————————")
	fmt.Printf("📅 开始执行每日任务 @ %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("——————————————")

	// 获取今日免费小说列表 📚
	fmt.Println("🔍 正在获取今日免费小说列表...")
	bookList, err := api.GetBooksList()
	if err != nil {
		fmt.Println("❌ 获取小说列表失败:", err)
		return false, fmt.Errorf("获取小说列表失败: %v", err)
	}

	fmt.Printf("✅ 成功获取%d本免费小说\n", len(bookList.Data.Data))
	if len(bookList.Data.Data) == 0 {
		fmt.Println("💤 今日没有免费小说，任务结束")
		return true, nil
	}

	// 处理每本小说 📖
	for i, book := range bookList.Data.Data {
		fmt.Printf("\n📚 处理第%d本小说: 《%s》\n", i+1, book.NovelName)
		fmt.Printf("🆔 小说ID: %s\n", book.NovelID)

		// 创建数据目录 📁
		dataDir := "data"
		_, err = os.Stat(dataDir)
		if os.IsNotExist(err) {
			fmt.Printf("📁 创建数据目录: %s\n", dataDir)
			err = os.Mkdir(dataDir, 0755)
			if err != nil {
				fmt.Println("❌ 创建数据目录失败:", err)
				return false, fmt.Errorf("创建数据目录失败: %v", err)
			}
			fmt.Println("✅ 数据目录创建成功")
		}

		// 创建小说文件 📄
		bookDir := dataDir + "/" + book.NovelName + ".txt"
		_, err = os.Stat(bookDir)
		if os.IsNotExist(err) {
			fmt.Printf("🆕 创建新小说文件: %s\n", bookDir)
			file, err := os.Create(bookDir)
			if err != nil {
				fmt.Println("❌ 创建小说文件失败:", err)
				return false, fmt.Errorf("创建小说文件失败: %v", err)
			}
			_ = file.Close()
			fmt.Println("✅ 小说文件创建成功")
		} else {
			fmt.Println("📝 小说文件已存在,跳过")
			continue
		}

		// 获取章节列表 📑
		fmt.Printf("🔍 获取《%s》的章节列表...\n", book.NovelName)
		chapterList, err := api.GetChapterList(book.NovelID)
		if err != nil {
			fmt.Println("❌ 获取章节列表失败:", err)
			return false, fmt.Errorf("获取章节列表失败: %v", err)
		}

		fmt.Printf("✅ 共获取%d个章节\n", len(chapterList.Chapterlist))
		var content string

		for j, chapter := range chapterList.Chapterlist {
			fmt.Printf("   📖 处理第%d章: %s (VIP: %v)\n", j+1, chapter.ChapterName, chapter.IsVip != 0)

			var chapterContent api.ChapterDetail
			if chapter.IsVip == 0 {
				fmt.Printf("   🆓 获取免费章节内容...\n")
				chapterContent, err = api.GetChapterContent(book.NovelID, chapter.ChapterID)
			} else {
				fmt.Printf("   💎 获取VIP章节内容...\n")
				chapterContent, err = api.GetVIPChapterContent(config.Token, book.NovelID, chapter.ChapterID)
			}

			if err != nil {
				fmt.Printf("   ❌ 获取章节内容失败: %s\n", err)
				return false, fmt.Errorf("获取章节内容失败: %v", err)
			}

			content += "第" + chapterContent.ChapterID + "章 " + chapterContent.ChapterName + "\n" + chapterContent.Content + "\n\n"
			fmt.Printf("   ✅ 第%d章处理完成\n", j+1)

			duration := time.Duration(config.Intervals.Chapter) * time.Millisecond
			fmt.Printf("   ⏸️ 休眠 %s 避免频繁请求...\n", duration)
			time.Sleep(duration)
		}

		// 写入文件 💾
		fmt.Printf("💾 正在将内容写入文件: %s\n", bookDir)
		err = os.WriteFile(bookDir, []byte(content), 0644)
		if err != nil {
			fmt.Println("❌ 写入文件失败:", err)
			return false, fmt.Errorf("写入文件失败: %v", err)
		}

		fmt.Printf("✅ 《%s》处理完成!\n", book.NovelName)

		duration := time.Duration(config.Intervals.Chapter) * time.Millisecond
		fmt.Printf("⏸️ 休眠 %s 避免频繁请求...\n", duration)
		time.Sleep(duration)
	}

	fmt.Println("——————————————")
	fmt.Printf("🎉 所有每日任务执行完成 @ %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("——————————————")
	return true, nil
}
