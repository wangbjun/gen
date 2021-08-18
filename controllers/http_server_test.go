package controllers

import (
	"gen/config"
	"gen/models"
	"gen/services"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 获取一个httpServer实例
func getHttpServer() *HTTPServer {
	cfg := config.NewConfig("../app.ini")

	err := cfg.Load()
	if err != nil {
		log.Fatalf("load ini file failed: %s", err)
	}
	cacheService := services.CacheService{
		Cfg: cfg,
	}
	err = cacheService.Init()
	if err != nil {
		log.Fatalf("cacheService init failed: %s", err)
	}
	sqlService := models.SQLService{Cfg: cfg}
	err = sqlService.Init()
	if err != nil {
		log.Fatalf("cacheService init failed: %s", err)
	}
	httpServer := &HTTPServer{
		Cfg:            cfg,
		ArticleService: &services.ArticleService{SQLStore: &sqlService, Cache: &cacheService},
		UserService:    &services.UserService{},
	}
	err = httpServer.Init()
	if err != nil {
		log.Fatalf("httpServer init failed: %s", err)
	}
	return httpServer
}

func TestHTTPServer(t *testing.T) {
	Convey("Test Server Status", t, func() {
		engine := getHttpServer()
		Convey("Status OK", func() {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/", nil)
			engine.ServeHTTP(recorder, request)

			So(recorder.Code, ShouldEqual, 200)
			So(recorder.Body.String(), ShouldContainSubstring, "Gen Web")
		})
		Convey("Status Not Found", func() {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/123456789", nil)
			engine.ServeHTTP(recorder, request)

			So(recorder.Code, ShouldEqual, 404)
			So(recorder.Body.String(), ShouldContainSubstring, "404")
		})
	})
}
