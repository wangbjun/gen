package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zht "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"strconv"
	"strings"
)

type errorCode int

const (
	Success      errorCode = 200
	Failed       errorCode = 500
	ParamError   errorCode = 400
	NotFound     errorCode = 404
	UnAuthorized errorCode = 401
)

var codeMsg = map[errorCode]string{
	Success:      "正常",
	Failed:       "系统异常",
	ParamError:   "参数错误",
	NotFound:     "记录不存在",
	UnAuthorized: "未授权",
}

type Controller struct {
	Page     int
	PageSize int
}

type Pagination struct {
	List       interface{} `json:"list"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
}

var BaseController = &Controller{}

var trans ut.Translator

// 注册validator中文翻译
func init() {
	uni := ut.New(en.New(), zh.New())
	trans, _ = uni.GetTranslator("zh")
	validate := binding.Validator.Engine().(*validator.Validate)
	_ = zht.RegisterDefaultTranslations(validate, trans)
}

func translate(err error) string {
	errors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	result := make([]string, 0)
	for _, err := range errors {
		errMsg := err.Translate(trans)
		if errMsg == "" {
			continue
		}
		result = append(result, errMsg)
	}
	return strings.Join(result, ",")
}

func (r *Controller) ParsePage(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 15
	}
	if pageSize > 1000 {
		pageSize = 1000
	}
	r.Page = page
	r.PageSize = pageSize
}

func (*Controller) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Gen Web")
}

func (*Controller) Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":     Success,
		"msg":      msg,
		"data":     data,
		"trace_id": ctx.GetString("trace_id"),
	})
}

func (*Controller) Failed(ctx *gin.Context, code errorCode, msg string) {
	errMsg := codeMsg[code] + ": " + msg
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code":     code,
		"msg":      errMsg,
		"data":     nil,
		"trace_id": ctx.GetString("trace_id"),
	})
	if code != Success {
		ctx.Set("error_code", int(code))
		ctx.Set("error_msg", msg)
	}
}
