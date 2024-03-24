package main

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"strings"
)

var log *zap.Logger

func init() {
	log, _ = zap.NewDevelopment()
}

// Fetch web content
func Fetch(url string) (content string, err error) {
	r, err := http.Get(url)
	if err != nil {
		log.Error("获取失败")
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	if r.StatusCode != http.StatusOK {
		log.Error(fmt.Sprintf("Error status code %v\n", r.StatusCode))
	}
	// 检查并转换
	bodyReader := bufio.NewReader(r.Body)
	e := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	body, err := io.ReadAll(utf8Reader)
	content = string(body)
	log.Info(fmt.Sprintf("一共有%v个链接", strings.Count(content, "<a")))
	return content, err
}

func DetermineEncoding(reader *bufio.Reader) encoding.Encoding {
	content, _ := reader.Peek(1024)
	e, _, _ := charset.DetermineEncoding(content, "")
	log.Info(fmt.Sprintf("Encoding: %v", e))
	return e
}

func main() {
	url := "https://www.thepaper.cn/"
	_, _ = Fetch(url)

}
