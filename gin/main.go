package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	engin := gin.Default() // Default 使用Logger和Recovery MiddleWare

	// 註冊validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("before_today", beforeToday)
	}

	engin.GET("/validation", func(ctx *gin.Context) {
		var user User
		if err := ctx.ShouldBind(&user); err != nil {
			msg := processErr(err)
			ctx.String(http.StatusBadRequest, "參數綁定失敗"+msg)
		} else {
			ctx.JSON(http.StatusOK, user)
		}
	})

	// 	Router分組
	{
		g1 := engin.Group("/v1")
		g1.Use(M6)
		g1.GET("/a", func(ctx *gin.Context) {
			ctx.String(200, "name=abc")
		})
		g1.GET("/b", func(ctx *gin.Context) {
			ctx.String(200, "age=18")
		})
	}

	{
		g2 := engin.Group("/v2")
		g2.GET("/a", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"name": "abc"})
		})
		g2.GET("/b", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"age": 18})
		})
	}

	engin.Use(M6)
	engin.GET("/1", M1, M2, M3, M4, M5)
	engin.GET("/2", M2, M3, M2, M4, M5)

	engin.GET("/home", gin.Logger(), homeHandler)

	url(engin)
	restful(engin)
	formBind(engin)
	jsonBind(engin)
	xmlBind(engin)
	yamlBind(engin)
	uriBind(engin)
	multiBind(engin)
	postForm(engin)
	postJson(engin)
	uploadFile(engin)
	uploadFiles(engin)

	if err := engin.Run("127.0.0.1:5678"); err != nil {
		panic(err)
	}
}

func M6(ctx *gin.Context) {
	slog.Info("visit", "path", ctx.Request.URL)
}

func M1(ctx *gin.Context) {
	ctx.String(200, "M1 Begin\n")
	ctx.Next()
	ctx.String(200, "M1 End\n")
}

func M2(ctx *gin.Context) {
	ctx.String(200, "Here is M2\n")
	// 就算沒有調用ctx.Next(), 也會進入下一個MiddleWare
}

func M3(ctx *gin.Context) {
	ctx.String(200, "Heere is M3\n")
}

func M4(ctx *gin.Context) {
	ctx.String(200, "Here is M4 Abort\n")
	ctx.Abort()
}

func M5(ctx *gin.Context) {
	ctx.String(200, "Here is M5\n")
}

func homeHandler(ctx *gin.Context) {
	fmt.Println("請求頭: ")
	for k, v := range ctx.Request.Header {
		fmt.Printf("%s=%s\n", k, v[0])
	}

	fmt.Println("請求體: ")
	io.Copy(os.Stdout, ctx.Request.Body)

	// 必須先設置響應頭
	ctx.Writer.Header().Add("language", "go")
	ctx.Header("Strict-Transport-Security", "max-age=31536000;includeSubDomains; preload")
	// 再設置響應碼
	ctx.Writer.WriteHeader(http.StatusOK)
	// 最後設置響應體
	ctx.Writer.WriteString("welcome")

	// Gin封裝
	ctx.String(200, "to home")
	ctx.JSON(200, map[string]any{"國文": 66, "數學": 77})
	ctx.JSON(200, gin.H{"物理": 35, "化學": 78})
}

// go build -tags=jsoniter -o gin.exe ./gin
