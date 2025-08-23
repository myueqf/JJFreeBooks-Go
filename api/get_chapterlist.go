package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChapterListRoot struct {
	Chapterlist         []Chapter `json:"chapterlist"`
	Count               int       `json:"count"`
	Isfree              int       `json:"isfree"`
	DaovLabelEndDate    string    `json:"daov_label_end_date"`
	VipChapterid        string    `json:"vipChapterid"`
	ProtectMeassge      string    `json:"protectMeassge"`
	BuyNoticeMeassge    string    `json:"buyNoticeMeassge"`
	BuyNoticeMeassge2   string    `json:"buy_notice_meassge"` // 注意：JSON 中有两个相似的字段名
	DiscountInfo        string    `json:"discount_info"`
	DiscountRatio       string    `json:"discount_ratio"`
	VipMonthFlag        string    `json:"vipMonthFlag"`
	HalfMoney           int       `json:"halfMoney"`
	HalfMoneyMessage    string    `json:"halfMoneyMessage"`
	LockInformation     string    `json:"lockInformation"`
	LockMessage         string    `json:"lockMessage"`
	Lockstatus          string    `json:"lockstatus"`
	BackBalance         int       `json:"backBalance"`
	BackRatio           int       `json:"backRatio"`
	Issign              bool      `json:"issign"`
	Editorbalance       string    `json:"editorbalance"`
	JjBalance           string    `json:"jjBalance"`
	ExtraChapterMeassge string    `json:"extraChapterMeassge"`
	VipShortFlag        string    `json:"vip_short_flag"`
	MonthRightEndTime   string    `json:"month_right_end_time"`
	MonthRightTip       string    `json:"month_right_tip"`
}

// Chapter 代表每一章的信息
type Chapter struct {
	NovelID              string   `json:"novelid"`
	ChapterID            string   `json:"chapterid"`
	ChapterType          string   `json:"chaptertype"`
	ChapterName          string   `json:"chaptername"`
	ChapterDate          string   `json:"chapterdate"`
	ChapterClick         string   `json:"chapterclick"` // 包含 "点击" 字样
	ChapterSize          string   `json:"chaptersize"`
	ChapterIntro         string   `json:"chapterintro"`
	IsLock               string   `json:"islock"`
	IsLockMessage        string   `json:"islockMessage"`
	IsVip                int      `json:"isvip"`
	DaovMsg              string   `json:"daov_msg"`
	Point                int      `json:"point"`
	OriginalPrice        int      `json:"originalPrice"`
	PointFreeVip         int      `json:"pointfreevip"`
	IsProtect            int      `json:"isProtect"`
	OriginalPriceMessage string   `json:"originalPriceMessage"`
	PointMeassge         string   `json:"pointMeassge"`
	ChapterMessage       string   `json:"chapterMessage"`
	LastPostTime         string   `json:"lastpost_time"`
	ExamineMessage       string   `json:"examineMessage"`
	IsEdit               string   `json:"isEdit"`
	Message              string   `json:"message"`
	Thank                int      `json:"thank"`
	TicketStarttime      string   `json:"ticketStarttime"`
	TicketEndtime        string   `json:"ticketEndtime"`
	Draft                Draft    `json:"draft"`
	ManageExplain        []string `json:"manageExplain"` // 空数组
	CheckChapterDate     int      `json:"checkChapterDate"`
	ExtraChapterType     string   `json:"extra_chapter_type"`
	ChapterRight         string   `json:"chapter_right"`
	ExtraChapterTip      string   `json:"extra_chapter_tip"`
}

// Draft 代表草稿信息
type Draft struct {
	Status  string `json:"status"`
	Color   string `json:"color"`
	Explain string `json:"explain"`
}

func GetChapterList(novelId string) ChapterListRoot {
	appUrl := fmt.Sprintf("https://app-cdn.jjwxc.com/androidapi/chapterList?novelId=%s", novelId)
	res, err := http.Get(appUrl)
	if err != nil {
		return ChapterListRoot{}
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	var result ChapterListRoot
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ChapterListRoot{}
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ChapterListRoot{}
	}
	return result
}
