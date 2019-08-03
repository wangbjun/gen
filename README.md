# Gen Web

## 介绍
基于Gin微框架封装的MVC Web框架demo，方便快速开发，主要使用以下开源组件：

* [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [github.com/go-redis/redis](https://github.com/go-redis/redis)
* [github.com/jinzhu/gorm](https://github.com/jinzhu/gorm)
* [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)
* [github.com/fvbock/endless](https://github.com/fvbock/endless)

目录结构清晰明了，支持平滑重启，demo是一个包含用户注册登录、发布文章、文章评论等功能的 restful api 应用

## 使用
项目使用 go module，建议git clone 到本地执行```go mod download```下载相关依赖，然后执行 storage/databases 目录下的sql，然后打开.env文件配置好数据库即可！
