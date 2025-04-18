package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Name       string    `form:"name" binding:"required"`
	Score      int       `form:"score" binding:"gt=0,required"`
	Enrollment time.Time `form:"enrollment" binding:"required,before_today" time_format:"2006-01-02" time_utc:"8"`
	Graduation time.Time `form:"graduation" binding:"required,gtfield=Enrollment" time_format:"2006-01-02" time_utc:"8"`
}

// 自定義validator
var beforeToday validator.Func = func(fl validator.FieldLevel) bool {
	// 透過類型斷言判斷資料類型
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		if date.Before(today) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func processErr(err error) string {
	if err == nil {
		return ""
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		msgs := make([]string, 0, 3)
		for _, validationErr := range validationErrs {
			msgs = append(msgs, fmt.Sprintf("字段 [%s] 不滿足條件[%s]", validationErr.Field(), validationErr.Tag()))
		}
		return strings.Join(msgs, ";")
	} else {
		return "invalid error"
	}
}
