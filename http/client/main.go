package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	myhttp "github.com/ohno104dev/go-web-framework/http/util"
)

func main() {
	// 	HttpObservation()
	// 	Get()
	Steam()
	// Student()
	// Head()
	// Post()
	// Cookie()
}

func Cookie() {
	fmt.Println(strings.Repeat("*", 30) + "COOKIE" + strings.Repeat("*", 30))

	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:5678/cookie", nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (x64)") // 偽造User-Agent
	request.Header.Add("user-role", "vvip")               // key會轉為規範大小寫

	// 添加多個Cookie
	request.AddCookie(
		&http.Cookie{
			Name:   "auth",
			Value:  "pass",
			Domain: "localhost",
			Path:   "/",
		},
	)

	// 所有的cookie都會放到一個http request header中, Cookie: [auth=pass; money=100]
	request.AddCookie(&http.Cookie{
		Name:  "money",
		Value: "100",
	})

	// cookie的key可以重複
	request.AddCookie(&http.Cookie{
		Name:  "money",
		Value: "800",
	})

	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	if resp, err := client.Do(request); err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		if values, exists := resp.Header["Set-Cookie"]; exists {
			fmt.Println(values[0])
			cookie, _ := http.ParseSetCookie(values[0])
			fmt.Println("Name:", cookie.Name)
			fmt.Println("Value:", cookie.Value)
			fmt.Println("Domain:", cookie.Domain)
			fmt.Println("MaxAge:", cookie.MaxAge)
			fmt.Println(strings.Repeat("-", 50))
		}
		os.Stdout.WriteString("\n\n")
	}
}

func Post() {
	fmt.Println(strings.Repeat("*", 30) + "POST" + strings.Repeat("*", 30))

	if resp, err := http.Post("http://127.0.0.1:5678/post", "text/plain", strings.NewReader("Hello Server")); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		if body, err := io.ReadAll(resp.Body); err == nil {
			fmt.Print(string(body))
		}
		os.Stdout.WriteString("\n\n")
	}

	bs, _ := json.Marshal(map[string]string{"name": "小安", "age": "18"})
	if resp, err := http.Post("http://127.0.0.1:5678/post", "application/json", bytes.NewReader(bs)); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}

	// PostForm()自動Content-Type設為application/x-www-form-urlencoded
	// 並把url.Values轉為URL-encode參數格式放到Request Body
	if resp, err := http.PostForm("http://127.0.0.1:5678/post", url.Values{"name": []string{"小陳"}, "age": []string{"33"}}); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}
}

// 可以用來檢查url是否存活, 只會取得HEAD部分, 不能取得response body
func Head() {
	fmt.Println(strings.Repeat("*", 30) + "HEAD" + strings.Repeat("*", 30))

	resp, err := http.Head("http://127.0.0.1:5678/get")
	if err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		if body, err := io.ReadAll(resp.Body); err == nil {
			fmt.Print(string(body))
		}
		os.Stdout.WriteString("\n\n")
	}
}

func Student() {
	fmt.Println(strings.Repeat("*", 30) + "GET" + strings.Repeat("*", 30))
	if resp, err := http.Get("http://127.0.0.1:5678/student"); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
		os.Stdout.WriteString("\n\n")
	}
}

// resp.Body為stream, 作用是避免一次性加載大量內容, 透過keep-alive做分段傳輸
func Steam() {
	fmt.Println(strings.Repeat("*", 30) + "GET HUGE BODY FOR STEAM TEST" + strings.Repeat("*", 30))
	if resp, err := http.Get("http://127.0.0.1:5678/stream"); err != nil {
		panic(err)
	} else {
		headerKey := http.CanonicalHeaderKey("cOnTeNt-LeNg") // 替換格式
		if h, ok := resp.Header[headerKey]; ok {
			if size, err := strconv.Atoi(h[0]); err == nil {
				haveRead := 0
				reader := bufio.NewReader(resp.Body)
				for {
					if bs, err := reader.ReadBytes('\n'); err == nil {
						haveRead += len(bs)
						progress := float64(haveRead) / float64(size)
						fmt.Printf("進度 %.2f%%, 內容 %s", 100*progress, string(bs))

						// 中途中止
						// if progress >= 0.5 {
						// 	resp.Body.Close()
						// 	return
						// }
					} else {
						if err == io.EOF {
							if len(bs) > 0 {
								fmt.Print(string(bs)) // 讀出剩餘的
							}
							break
						} else {
							fmt.Printf("read response body error: %s\n", err)
						}
					}
				}
				resp.Body.Close()
			}
		} else {
			println("xxxxx")
		}

	}
}

func HttpObservation() {
	fmt.Println(strings.Repeat("*", 30) + "GET" + strings.Repeat("*", 30))
	resp, err := http.Get("http://127.0.0.1:5678/obs?name=abc&age=20")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("response proto: %s\n", resp.Proto)
	if major, minor, ok := http.ParseHTTPVersion(resp.Proto); ok {
		fmt.Printf("http major version %d, http minor version %d\n", major, minor)
	}

	fmt.Printf("response status: %s\n", resp.Status)
	fmt.Printf("response status code: %d\n", resp.StatusCode)

	for k, v := range resp.Header {
		fmt.Printf("%s: %v\n", k, v)
	}

	fmt.Println("response body:")
	io.Copy(os.Stdout, resp.Body)
	os.Stdout.WriteString("\n\n")

}

func Get() {
	fmt.Println(strings.Repeat("*", 30) + "GET" + strings.Repeat("*", 30))

	resp, err := http.Get("http://127.0.0.1:5678/get?" + myhttp.EncodeUrlParams(map[string]string{"name": "你好啊 &&& %%% 這是測試!!!", "age": "20"}))
	if err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("response status: %s\n", resp.Status)
		fmt.Println("response body:")
		if body, err := io.ReadAll(resp.Body); err == nil {
			fmt.Print(string(body))
		}
		os.Stdout.WriteString("\n\n")
	}
}
