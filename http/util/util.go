package util

import (
	"net/url"
	"strings"
)

// 透過QueryEscape 安全轉譯漢字, 空格, 特殊符號
func EncodeUrlParams(params map[string]string) string {
	sb := strings.Builder{}
	for k, v := range params {
		sb.WriteString(url.QueryEscape(k))
		sb.WriteString("=")
		sb.WriteString(url.QueryEscape(v))
		sb.WriteString("&")
	}

	return sb.String()
}

// 透過QueryEscape 安全解譯漢字, 空格, 特殊符號
func ParseUrlParams(rawQuery string) map[string]string {
	params := make(map[string]string, 10)
	args := strings.Split(rawQuery, "&")

	for _, ele := range args {
		arr := strings.Split(ele, "=")
		if len(arr) == 2 {
			k, _ := url.QueryUnescape(arr[0])
			v, _ := url.QueryUnescape(arr[1])
			params[k] = v
		}
	}

	return params
}
