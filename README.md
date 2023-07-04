# Gen Web

## 介绍
一个基于Gin框架封装的脚手架工具，开箱即用，便于用Go快速开发一些Web API，使用以下开源组件：

* [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [gorm.io/gorm](https://gorm.io/gorm)
* [go.uber.org/zap](https://go.uber.org/zap)
* [github.com/go-playground/validator/v10](https://github.com/go-playground/validator/v10)
* [github.com/go-redis/redis](https://github.com/go-redis/redis)

项目目录结构清晰明了，简单易用，快速上手，包含了一个用户注册、登录、文章增删改查等功能的 Restful API 应用，仅供参考！

主要包含以下API：

|METHOD|URI|DESCRIPTION|
|---|---|---|
|GET|/|默认首页
|POST|/api/v1/user/register|用户注册
|POST|/api/v1/user/login|用户登录
|POST|/api/v1/user/logout|用户登出
|GET|/api/v1/articles|文章列表
|POST|/api/v1/articles|发布文章
|GET|/api/v1/articles/:id|文章详情
|PUT|/api/v1/articles/:id|修改文章
|DELETE|/api/v1/articles/:id|删除文章
|POST|/api/v1/articles/:id/comments|添加文章评论

## 架构
之前有一版是借鉴了著名Go开源项目 [Grafana](https://github.com/grafana/grafana) 的设计，使用了依赖注入机制，但是感觉过于复杂，不容易理解和使用，所以又改了。

本着简单易用易修改的原则，采用了包全局变量的方式初始化配置、日志、DB连接，项目目录如下：
```
├── config //配置
├── controllers //控制器
├── log //日志
├── middleware //中间件
├── models //数据表模型
├── router //路由
├── services //服务
└── utils //工具函数
```
在 ```main.go``` 里面依次初始化各个组件，清晰明了：
```go
func main() {
    var configFile string
    flag.StringVar(&configFile, "conf", "app.ini", "config file path")
    flag.Parse()

    // 加载配置
    cfg := config.InitConfig(configFile)
    err := cfg.Load()
    if err != nil {
        panic(fmt.Sprintf("load config failed, file: %s, error: %s", configFile, err))
    }

    // 初始化日志
    log.Init(cfg)
    defer func() {
        if err := log.Logger.Sync(); err != nil {
            fmt.Printf("Failed to close log: %s\n", err)
        }
    }()

    // 初始化数据库
    err = models.InitDB(cfg)
    if err != nil {
        panic(fmt.Sprintf("init db failed, error: %s", err))
    }

    // 启动Web服务
    err = startServer(cfg)
    if err != nil {
        panic(fmt.Sprintf("Server started failed: %s", err))
    }
}
```

## 代码介绍
在services文件夹下包含了一些服务的代码文件。

项目整体是一个3层架构，即控制器层、Service层、模型层。

个人理解，控制器层主要做一些接口参数校验等工作，模型层主要是数据操作，Service层才是主要的业务逻辑。

数据库相关配置在models/db.go里面，也是一个服务，主要作用是根据配置，初始化数据库连接，支持多数据库配置、支持Sql日志记录。

项目使用了Gorm（2.0版本），具体详细用法可以参考官方文档。

config/config.go是配置文件的一些加载逻辑，可以根据自己需求适当的修改优化。

关于接口参数，建议POST、PUT统一使用JSON形式，在模型层里面定义好相应的结构体，参数的校验采用了```go-playground/validator/v10```库，直接在结构体Tag里面标记即可，详细用法请参考其官方文档。
```go
type CreateArticleCommand struct {
    Id      int
    UserId  int
    Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
    Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}

type UpdateArticleCommand struct {
    Id      int
    UserId  int
    Title   string `form:"title" json:"title" binding:"gt=1,lt=100"`
    Content string `form:"content" json:"content" binding:"gt=1,lt=2000"`
}
```
## 命名规范
个人建议参考以下规范：

- 文件夹名全部小写，多个单词的话直接相连。如 ```eventhandler```
- 文件名小写下划线，如 ```article_controller.go```
- 变量名驼峰，如 ```var userId int```，```pageSize := 10```，首字母是否大写根据实际需要（是否要公开）
- 结构体名驼峰，如 ```type userService struct ```，首字母是否大写根据实际需要（是否要公开）
- 函数名驼峰，如果 ```func getUserById()```，首字母是否大写根据实际需要（是否要公开）

其它规范建议以Goland内置标准为准，一般情况下IDE都会提示，建议遵循。

## 使用
建议直接clone本项目，然后删除多余的控制器、模型等文件，根据自己需求调整即可，Golang的项目真滴很简单！