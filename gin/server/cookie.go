package main

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	counter sync.Map
)

func init() {
	counter.Store("idgnoiewfq", 0)
}

// 紀錄每個token調用接口的次數
func CountMD(ctx *gin.Context) {
	token, err := ctx.Cookie("token")
	if err != nil {
		ctx.Abort()
	}

	if v, exists := counter.Load(token); !exists {
		ctx.Abort()
	} else {
		c, _ := v.(int)
		counter.Store(token, c+1)
		fmt.Println("visit counter, token=", token, "count=", c+1)
	}
}

func SetCookie(ctx *gin.Context) {
	name := "language"
	value := "go"
	maxAge := 86400 * 7 // 如果不設置過期時間, 默認關閉瀏覽器後cookie被刪除
	path := "/"         // cookie存放的目錄

	// cookie從屬的domain名
	// 如果不指定默認為host
	// 如果指定的domain是一級domain名, 則以下的二級domain名也可以訪問
	domain := "www.google.com"

	secure := false  // 是否只通過https訪問
	httpOnly := true // 是否允許通過js獲取cookie, 防止XSS攻擊

	// SetCookie只能執行一次
	// 對應的response header key 是"Set-Cookie"
	ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}
