package main

import (
	"net/http"

	"github.com/ohno104dev/go-web-framework/gin/idl"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func text(engine *gin.Engine) {
	engine.GET("/user/text", func(ctx *gin.Context) {
		// response Content-Type: text/plain
		ctx.String(http.StatusOK, "hi boy")
	})
}

func json0(engine *gin.Engine) {
	engine.GET("/user/json0", func(ctx *gin.Context) {
		// response Content-Type: application/json
		var stu struct {
			Name string `json:"name"`
			Addr string `json:"addr"`
		}

		s, _ := sonic.MarshalString(stu)
		ctx.Request.Header.Add("Content-Type", "application/json")
		ctx.String(http.StatusOK, s)
	})
}

func json1(engine *gin.Engine) {
	engine.GET("/user/json1", func(ctx *gin.Context) {
		// response Content-type: application/json
		ctx.JSON(http.StatusOK, gin.H{
			"name": "zcy",
			"addr": "BK",
		})
	})
}

func json2(engine *gin.Engine) {
	var stu struct {
		Name string
		Addr string
	}
	stu.Name = "abc"
	stu.Addr = "NY"
	engine.GET("/user/json2", func(ctx *gin.Context) {
		// response Content-type: application/json
		ctx.JSON(http.StatusOK, stu)
	})
}

// 使用 JSONP 可以向不同domain的服務器請求數據
func jsonp(engine *gin.Engine) {
	var stu struct {
		Name string
		Addr string
	}

	stu.Name = "ccc"
	stu.Addr = "KU"
	engine.GET("/user/jsonp", func(ctx *gin.Context) {
		// 如果請求參數裡有callback=xxx, 則response Content-Type為: pplication/javascript
		// 否則response Content-Type: application/json
		ctx.JSONP(http.StatusOK, stu)
	})
}

func xml(engine *gin.Engine) {
	type Student struct {
		Name string
		Addr string
	}

	var stu Student
	stu.Name = "qqq"
	stu.Addr = "TY"
	engine.GET("/user/xml", func(ctx *gin.Context) {
		// response Content-Type: application/xml
		ctx.XML(http.StatusOK, stu)
	})
}

func protoBuf(engine *gin.Engine) {
	stu := &idl.Student{
		Name:    "zcy",
		Address: "BJ",
	}

	engine.GET("/user/pb", func(ctx *gin.Context) {
		// response Content-Type: application/x-protobuf
		ctx.ProtoBuf(http.StatusOK, stu)
	})
}

func html(engine *gin.Engine) {
	// engine.LoadHTMLGlob("web/static/")
	engine.LoadHTMLFiles("gin/static/template.html")
	// Load後, 可以直接用template.html
	engine.GET("/user/html", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "template.html", gin.H{
			"title": "用戶訊息",
			"name":  "zcx",
			"addr":  "AC",
		})
	})
}

func redirect(engine *gin.Engine) {
	engine.GET("/user/old_page", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "http://localhost:5678/user/html")
	})
}
