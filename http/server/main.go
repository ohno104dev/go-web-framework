package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	myhttp "github.com/ohno104dev/go-web-framework/http/util"
)

func main() {
	http.HandleFunc("/obs", HttpObservation)
	http.HandleFunc("/get", Get)
	http.HandleFunc("/stream", StreamBody)
	http.HandleFunc("/student", Student)
	http.HandleFunc("/post", Post)
	http.HandleFunc("/cookie", Cookie)

	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		panic(err)
	}
}

func Cookie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request header:")
	for key, value := range r.Header {
		fmt.Println(key, value)
	}

	// 可以直接使用r.Cookies()
	if values, exists := r.Header["Cookie"]; exists {
		cookies, _ := http.ParseCookie(values[0])
		for _, cookie := range cookies {
			fmt.Printf("%s: %s\n", cookie.Name, cookie.Value)
		}
		fmt.Println(strings.Repeat("*", 60))
	}

	// Set-Cookie, 服務端只能返回一個
	expiration := time.Now().Add(30 * 24 * time.Hour)
	cookie1 := http.Cookie{Name: "csrftoken", Value: "abcdabcd", Expires: expiration, Domain: "localhost", Path: "/"}
	cookie2 := http.Cookie{Name: "jwt", Value: "77777", Expires: expiration, Domain: "localhost", Path: "/"}
	http.SetCookie(w, &cookie1) // SetCookie只能執行一次
	http.SetCookie(w, &cookie2) // 無效
}

func Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if ct, exists := r.Header["Content-Type"]; exists {
		switch ct[0] {
		case "text/plain":
			io.Copy(w, r.Body)
		case "application/json":
			body, err := io.ReadAll(r.Body)
			if err == nil {
				params := make(map[string]string, 0)
				if err := json.Unmarshal(body, &params); err == nil {
					fmt.Fprintf(w, "your name is %s, age is %s \n", params["name"], params["age"])
				}
			} else {
				fmt.Println("read request body error", err)
				w.WriteHeader(http.StatusBadRequest)
			}
		case "application/x-www-form-urlencoded":
			body, err := io.ReadAll(r.Body)
			if err == nil {
				params := myhttp.ParseUrlParams(string(body))
				fmt.Println("request body", string(body))
				fmt.Fprintf(w, "your name is %s, age is %s \n", params["name"], params["age"])
			} else {
				fmt.Println("read request body error", err)
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}

func Student(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type Student struct {
		Id     int
		Name   string
		Gender string
		Score  int
	}

	student := []Student{
		{1, "張三", "男", 80},
		{2, "李四", "女", 77},
	}

	tmpl, err := template.ParseFiles("./http/server/student.html")
	if err != nil {
		fmt.Println("create template failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, student)
}

func StreamBody(w http.ResponseWriter, r *http.Request) {
	line := []byte("This is a mock test for send huge body data. Heavy is the head who wears the crown. \n")
	// const rp = 1000_000_000
	const rp = 10
	totalSize := rp * len(line)
	hkey := http.CanonicalHeaderKey("cOntEnt-lEnG")
	w.Header().Add(hkey, strconv.Itoa(totalSize))
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "不支持Flush", http.StatusInternalServerError)
		return
	}

	for i := range rp {
		// 即使不顯式調用Flush(), Write()內容足夠大時也會觸發Flush()
		if _, err := w.Write(line); err != nil {
			fmt.Printf("rp: [%d] send error: %s\n", i, err)
			break
		} else {
			flusher.Flush()
			time.Sleep(time.Second * 1)
		}
	}
	fmt.Println(strings.Repeat("*", 60))
}

func Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	params := myhttp.ParseUrlParams(r.URL.RawQuery)
	fmt.Fprintf(w, "your name is %s, age is %s\n", params["name"], params["age"])
	fmt.Println(strings.Repeat("*", 60))
}

func HttpObservation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Printf("request method: %s\n", r.Method)
	fmt.Printf("request host: %s\n", r.Host)
	fmt.Printf("request url: %s\n", r.URL)
	fmt.Printf("request proto: %s\n", r.Proto)
	fmt.Println("request header")
	for k, v := range r.Header {
		fmt.Printf("%s: %v\n", k, v)
	}
	fmt.Println()
	fmt.Printf("request body: ")
	// io.Copy(os.Stdout, r.Body)
	if body, err := io.ReadAll(r.Body); err == nil {
		fmt.Println(string(body))
	}
	fmt.Println()
	r.Body.Close()

	// 注意順序 Header -> WriteHeader
	w.Header().Add("tRAce-id", "198345906")
	w.WriteHeader(http.StatusBadRequest) // WriteHeader之後會lock Header, 之後無法修改Header
	w.Write([]byte("Hello World\n"))
	w.Write([]byte("Hello Boy\n"))
	fmt.Fprint(w, "Hello Girl\n")
	fmt.Println(strings.Repeat("*", 60))
}
