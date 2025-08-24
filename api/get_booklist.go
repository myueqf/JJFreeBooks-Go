package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// BookListRoot 是整个 JSON 响应的根结构
type BookListRoot struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

// Data 包含主要的业务数据
type Data struct {
	ChannelID       string      `json:"channelid"`
	NaturalRankID   int         `json:"natural_rank_id"`
	Data            []NovelData `json:"data"`
	ChannelRule     string      `json:"channelRule"`
	ChannelRuleDown string      `json:"channelRuleDown"`
}

// NovelData 代表每本小说的数据
type NovelData struct {
	NovelID         string `json:"novelId"`
	NovelName       string `json:"novelName"`
	AuthorID        string `json:"authorId"`
	AuthorName      string `json:"authorName"`
	Cover           string `json:"cover"`
	Local           string `json:"local"`
	LocalImg        string `json:"localImg"`
	NovelIntroshort string `json:"novelIntroshort"`
	NovelIntroShort string `json:"novelIntroShort"` // 注意：JSON 中有两个相似的字段
	NovelIntro      string `json:"novelIntro"`
	NovelStep       string `json:"novelStep"`
	Tags            string `json:"tags"`
	FreeDate        string `json:"freeDate"`
	NowFree         string `json:"nowFree"`
	NovelSize       string `json:"novelSize"`
	NovelSizeformat string `json:"novelSizeformat"`
	NovelClass      string `json:"novelClass"`
	IsVipMonth      string `json:"isVipMonth"`
	RecommendInfo   string `json:"recommendInfo"` // 这个字段本身是 JSON 字符串
}

func GetBooksList() (BookListRoot, error) {
	now := time.Now()
	date := now.Format("2006-01-02")

	channelBody := fmt.Sprintf(`{"date_free_%s":{"offset":"0","limit":"40"}}`, date)

	escapedChannelBody := url.QueryEscape(channelBody)

	appUrl := fmt.Sprintf("https://app-cdn.jjwxc.com/bookstore/getFullPageV1?channelBody=%s&channelMore=1", escapedChannelBody)

	res, err := http.Get(appUrl)
	if err != nil {
		return BookListRoot{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	var result BookListRoot
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return BookListRoot{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return BookListRoot{}, err
	}
	return result, nil
}
