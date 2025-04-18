package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/ohno104dev/go-web-framework/gin/idl"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 從GET請求的URL中獲取參數
func url(engine *gin.Engine) {
	// http://127.0.0.1:5678/student?name=abc&
	engine.GET("/student", func(ctx *gin.Context) {
		a := ctx.Query("name")
		b := ctx.DefaultQuery("addr", "tw")
		ctx.String(http.StatusOK, a+" live in "+b)
	})
}

// 從Restful風格的URL中獲取參數
func restful(engine *gin.Engine) {
	// http://127.0.0.1:5678/student/aaa/us
	engine.GET("/student/:name/*addr", func(ctx *gin.Context) {
		name := ctx.Param("name")
		addr := ctx.Param("addr") //*多級對應, 會包含'/us'
		ctx.String(http.StatusOK, name+" live in "+addr)
	})
}

func postForm(engine *gin.Engine) {
	engine.POST("/student/form", func(ctx *gin.Context) {
		name := ctx.PostForm("username")
		addr := ctx.DefaultPostForm("addr", "Japan")
		ctx.String(http.StatusOK, name+" live in "+addr)
	})
}

type Student struct {
	Name     string   `form:"username" json:"name" uri:"user" xml:"user" yaml:"user" binding:"required"`
	Addr     string   `form:"addr" json:"addr" uri:"addr" xml:"addr" yaml:"addr" binding:"required"`
	Keywords []string `form:"keywords"`
}

// *** ctx.ShouldBindBodyWith 會在綁定之前將body儲存到Context, 可以重複讀取body, 但會對性能造成輕微影響
func multiBind(engine *gin.Engine) {
	engine.POST("/stu/multi_type", func(ctx *gin.Context) {
		var stu Student
		var stu2 idl.Student
		if err := ctx.ShouldBindBodyWith(&stu, binding.JSON); err == nil {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		} else if err := ctx.ShouldBindBodyWith(&stu, binding.XML); err == nil {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		} else if err := ctx.ShouldBindBodyWith(&stu, binding.YAML); err == nil {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		} else if err := ctx.ShouldBindBodyWith(&stu2, binding.ProtoBuf); err == nil {
			ctx.String(http.StatusOK, stu2.Name+" live in "+stu2.Address)
		} else {
			ctx.String(http.StatusBadRequest, "不支持的參數類型")
		}
	})
}

func uriBind(engine *gin.Engine) {
	//Get 請求的參數在uri裡
	engine.GET("stu/uri/:user/*addr", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBindUri(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse parameter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func xmlBind(engine *gin.Engine) {
	engine.POST("/stu/xml", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBindXML(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse parameter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func yamlBind(engin *gin.Engine) {
	engin.POST("/stu/yaml", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBindYAML(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse parameter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func jsonBind(engine *gin.Engine) {
	engine.POST("stu/json", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBindJSON(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse parameter filed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func formBind(engine *gin.Engine) {
	engine.POST("/stu/form", func(ctx *gin.Context) {
		var stu Student
		if err := ctx.ShouldBind(&stu); err != nil {
			fmt.Println(err)
			ctx.String(http.StatusBadRequest, "parse parameter failed")
		} else {
			ctx.String(http.StatusOK, stu.Name+" live in "+stu.Addr)
		}
	})
}

func postJson(engine *gin.Engine) {
	engine.POST("/student/json", func(ctx *gin.Context) {
		var stu Student
		bs, _ := io.ReadAll(ctx.Request.Body)
		if err := json.Unmarshal(bs, &stu); err == nil {
			name := stu.Name
			addr := stu.Addr
			ctx.String(http.StatusOK, name+" live in "+addr)
		}
	})
}

func uploadFile(engine *gin.Engine) {
	engine.MaxMultipartMemory = 8 << 20 //限制表單上傳大小為8M, Default是32M
	engine.POST("/upload", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			fmt.Printf("get file error %v\n", err)
			ctx.String(http.StatusInternalServerError, "upload file failed")
		} else {
			if err = ctx.SaveUploadedFile(file, "./data/"+file.Filename); err == nil {
				ctx.String(http.StatusOK, file.Filename)
			} else {
				fmt.Printf("save file to %s failed: %v\n", "./data"+file.Filename, err)
			}
		}
	})
}

func uploadFiles(engine *gin.Engine) {
	engine.POST("/uploads", func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
		} else {
			// 從MultipartForm中獲取上傳的文件
			files := form.File["files"]
			for _, file := range files {
				ctx.SaveUploadedFile(file, "./data/"+file.Filename)
			}

			ctx.String(http.StatusOK, "upload "+strconv.Itoa(len(files))+" files")
		}
	})
}
