package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	Get("/home")
}

func Get(path string) {
	fmt.Print("GET " + path + " ")
	resp, err := http.Get("http://127.0.0.1:5678" + path)
	if err != nil {
		panic(err)
	}
	processResponse(resp)
}

func processResponse(resp *http.Response) {
	defer resp.Body.Close()
	fmt.Println("響應頭: ")
	for k, v := range resp.Header {
		fmt.Printf("%s=%s\n", k, v[0])
	}
	fmt.Println("響應體: ")
	fmt.Println("response code: ", resp.StatusCode)
	io.Copy(os.Stdout, resp.Body)
	os.Stdout.WriteString("\n\n")
}
