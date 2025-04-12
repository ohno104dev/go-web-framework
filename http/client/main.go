package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	myhttp "github.com/ohno104dev/go-web-framework/http/util"
)

func main() {
	HttpObservation()
	Get()
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
