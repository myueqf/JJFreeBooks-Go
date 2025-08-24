package api

import (
	"JJFreeBooks/crypto"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ChapterDetail struct {
	ChapterID       string `json:"chapterId"`
	ChapterName     string `json:"chapterName"`
	ChapterIntro    string `json:"chapterIntro"`
	ChapterSize     string `json:"chapterSize"`
	ChapterDate     string `json:"chapterDate"`
	SayBody         string `json:"sayBody"`
	UpDown          int    `json:"upDown"`
	Update          int    `json:"update"`
	Content         string `json:"content"`
	IsVip           int    `json:"isvip"`
	AuthorID        string `json:"authorid"`
	AutoBuyStatus   string `json:"autobuystatus"`
	NoteIsLock      int    `json:"noteislock"`
	SayBodyV2       string `json:"sayBodyV2"`
	ShowSayBodyPage string `json:"show_saybody_page"`
}

func GetChapterContent(novelId, chapterId int) (ChapterDetail, error) {
	appUrl := fmt.Sprintf("https://app-cdn.jjwxc.com/androidapi/chapterContent?novelId=%s&chapterId=%s", strconv.Itoa(novelId), strconv.Itoa(chapterId))
	res, err := http.Get(appUrl)
	if err != nil {
		return ChapterDetail{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ChapterDetail{}, err
	}
	var result ChapterDetail
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ChapterDetail{}, err
	}
	return result, nil
}

func GetVIPChapterContent(token string, novelId, chapterId int) (ChapterDetail, error) {
	timestamp := time.Now().UnixMilli()
	ciphertextStr := fmt.Sprintf("%s:%s:%s:%s", strconv.Itoa(int(timestamp)), token, strconv.Itoa(novelId), strconv.Itoa(chapterId))

	ciphertext, err := crypto.DesEncrypt([]byte(ciphertextStr), []byte("KW8Dvm2N"), []byte("1ae2c94b"))
	if err != nil {
		return ChapterDetail{}, err
	}
	escapedChannelBody := url.QueryEscape(ciphertext)
	fmt.Println(escapedChannelBody)
	appUrl := fmt.Sprintf("https://android.jjwxc.net/androidapi/chapterContent?readState=readahead&versionCode=454&sign=%s", escapedChannelBody)

	req, err := http.NewRequest("GET", appUrl, nil)
	if err != nil {
		return ChapterDetail{}, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Linux; Android 15; PJX110 Build/UKQ1.231108.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/139.0.7258.94 Mobile Safari/537.36/JINJIANG-Android/454(PJX110;Scale/2.55;isHarmonyOS/false)")
	req.Header.Add("Accept-Encoding", "")
	req.Header.Add("local_maxchapterid", "71")
	req.Header.Add("cacheshowed", "false")
	req.Header.Add("referer", "http://android.jjwxc.net/?v=454")
	req.Header.Add("not_tip", "readahead")
	req.Header.Add("versiontype", "reading")
	req.Header.Add("versiontype", "reading")
	req.Header.Add("source", "android")
	req.Header.Add("versioncode", "454")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ChapterDetail{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ChapterDetail{}, err
	}
	var result ChapterDetail
	var dest string
	var keyIv = crypto.KeyIv{}
	var key, iv []byte
	accesskey := res.Header.Get("accesskey")
	keystring := res.Header.Get("keystring")

	if json.Valid(body) {
		var r map[string]interface{}
		err = json.Unmarshal(body, &r)
		if err != nil {
			return ChapterDetail{}, err
		}
		fmt.Println(r)
		if content, ok := r["content"].(string); ok {
			dest = content
		} else {
			return ChapterDetail{}, fmt.Errorf("字段 'content' 不存在或不是字符串")
		}
	} else {
		dest, keyIv, err = crypto.DynamicDecrypt(string(body), accesskey, keystring)
		fmt.Println(dest)
		if err != nil {
			return ChapterDetail{}, err
		}
	}
	key = keyIv.Key
	iv = keyIv.Iv
	fmt.Println(hex.EncodeToString(key))
	fmt.Println(hex.EncodeToString(iv))

	var data string
	data, err = crypto.DesDecrypt(dest, key, iv)
	if err != nil {
		return ChapterDetail{}, err
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return ChapterDetail{}, err
	}
	return result, nil
}
