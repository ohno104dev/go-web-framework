package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	engin := gin.Default()
	engin.GET("/home", homeHandler)
	if err := engin.Run("127.0.0.1:5678"); err != nil {
		panic(err)
	}
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
