package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/ohno104dev/go-web-framework/gin/idl"

	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

func main() {
	Get("/home")
	PostForm("/student/form", Student{Name: "小陳", Addr: "Us"})
	PostJson("/student/json", Student{Name: "小王", Addr: "Us"})
	PostXml("/stu/xml", Student{Name: "小黃", Addr: "Uk"})
	PostYaml("/stu/yaml", Student{Name: "小張", Addr: "Aus"})
	PostPb("/stu/multi_type", Student{Name: "小吳", Addr: "Cn"})
}

type Student struct {
	Name     string   `form:"username" json:"name" uri:"user" xml:"user" yaml:"user" binding:"required"`
	Addr     string   `form:"addr" json:"addr" uri:"addr" xml:"addr" yaml:"addr" binding:"required"`
	Keywords []string `form:"keywords"`
}

func PostPb(path string, stu Student) {
	fmt.Print("post pb " + path + " ")
	inst := idl.Student{Name: stu.Name, Address: stu.Addr}
	if bs, err := proto.Marshal(&inst); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "", bytes.NewReader(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("proto marshal failed", "error", err)
	}
}

func PostXml(path string, stu Student) {
	fmt.Print("post xml " + path + " ")
	if bs, err := xml.Marshal(stu); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "", bytes.NewReader(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("xml marshal failed", "error", err)
	}
}

func PostYaml(path string, stu Student) {
	fmt.Print("post yaml " + path + " ")
	if bs, err := yaml.Marshal(stu); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "", bytes.NewReader(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("yaml marshal failed", "error", err)
	}
}

func PostJson(path string, stu Student) {
	fmt.Print("post json " + path + " ")
	if bs, err := json.Marshal(stu); err == nil {
		if resp, err := http.Post("http://127.0.0.1:5678"+path, "", bytes.NewBuffer(bs)); err != nil {
			panic(err)
		} else {
			processResponse(resp)
		}
	} else {
		slog.Error("json marchal failed", "error", err)
	}
}

func PostForm(path string, stu Student) {
	fmt.Print("post form " + path + " ")
	// PostForm()會自動把請求頭的Content-Type設為Application/x-www-form-urlencoded, 並把url.Values轉為URL-encoded 參數格式放到請求體
	if resp, err := http.PostForm("http://127.0.0.1:5678"+path, url.Values{
		"username": {stu.Name}, "addr": {stu.Addr}}); err != nil {
		panic(err)
	} else {
		processResponse(resp)
	}

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
