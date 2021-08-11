package controllers

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleController_GetAll(t *testing.T) {
	Convey("Test GetAll", t, func() {
		engine := getHttpServer()
		Convey("Normal", func() {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/api/v1/articles/1", nil)
			engine.ServeHTTP(recorder, request)

			So(recorder.Code, ShouldEqual, 200)
			So(recorder.Body.String(), ShouldContainSubstring, "88888")
		})
		Convey("Not Found", func() {
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/api/v1/articles/99999", nil)
			engine.ServeHTTP(recorder, request)

			So(recorder.Code, ShouldEqual, 200)
			So(recorder.Body.String(), ShouldContainSubstring, "not found")
		})
	})
}
