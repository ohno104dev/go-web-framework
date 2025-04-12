package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	myhttp "github.com/ohno104dev/go-web-framework/http/util"
)

func main() {
	http.HandleFunc("/obs", HttpObservation)
	http.HandleFunc("/get", Get)

	if err := http.ListenAndServe("127.0.0.1:5678", nil); err != nil {
		panic(err)
	}
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
