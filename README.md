# Gen Web

## 介绍
基于Gin微框架封装的MVC Web框架demo，方便快速开发，主要使用以下开源组件：

* [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [github.com/go-redis/redis](https://github.com/go-redis/redis)
* [github.com/jinzhu/gorm](https://github.com/jinzhu/gorm)
* [github.com/fvbock/endless](https://github.com/fvbock/endless)

目录结构清晰明了，支持平滑重启，demo是一个包含用户注册登录、发布文章、文章评论等功能的 restful api 应用

主要包含以下API：

|METHOD|URI|DESCRIPTION|
|---|---|---|
|GET|/|默认首页
|GET|/api/v1/articles|文章列表
|POST|/api/v1/articles|发布文章
|GET|/api/v1/articles/:id|文章详情
|POST|/api/v1/articles/:id|修改文章
|DELETE|/api/v1/articles/:id|删除文章
|POST|/api/v1/articles/:id/comments|添加文章评论
|GET|/api/v1/articles/:id/comments|文章评论列表
|POST|/api/v1/user/register|用户注册
|POST|/api/v1/user/login|用户登录
|POST|/api/v1/user/logout|用户登出

## 使用
项目使用 go module，建议git clone 到本地执行```go mod download```下载相关依赖，然后执行 storage/databases 目录下的sql，然后打开 app.ini 文件配置好数据库即可！
