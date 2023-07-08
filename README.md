# Gen Web

## 介绍
众所周知，Go在Web领域有很多应用，不少公司拿Go来写一些接口，不仅在安全性可靠性上面有优势，而且得益于Go强大的协程机制，性能也很高，但是Go至今仍然缺乏一个全栈的Web框架。

或者说，并不是缺乏，因为也一直有各种Web框架诞生，但是至今仍然没有一个框架得到广泛的认同和普及，没有一个框架能达到Java界的Spring，或者PHP界的Laravel那种程度。

在我看来，一个Web框架至少应该包括以下几个方面：

- 路由
- 日志
- 数据库|ORM
- 配置管理
- 控制器
- 中间件
- 鉴权

其实这些组件单独拎出来都有很多知名的项目，但是组合到一起的并不多，比如Gin这个框架，非常优秀，但是缺少配置、数据库的组件，它只包括一个核心的Web模块，也有一些框架虽然包括了这些组件，但是过于复杂了，也许你并不需要这些功能。

这个项目就是一个基于Gin框架封装的脚手架工具，在Gin框架的基础上增加了日志、ORM、配置等模块，其实只是把一些开源的组件拼了一下，便于用Go快速开发一些Web API，使用以下开源组件：

* [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [gorm.io/gorm](https://gorm.io/gorm)
* [go.uber.org/zap](https://go.uber.org/zap)
* [github.com/go-playground/validator/v10](https://github.com/go-playground/validator/v10)
* [github.com/go-redis/redis](https://github.com/go-redis/redis)

项目主打的就是一个简单易用，快速上手，如果你觉得哪里不合适，自己改一改就行了，包含了一个文章增删改查等功能的 Restful API 应用，如果你喜欢Gin框架，不妨参考一下！

主要包含以下API：

| METHOD |URI|DESCRIPTION|
|--------|---|---|
| GET    |/|默认首页
| GET    |/api/v1/articles|文章列表
| POST   |/api/v1/articles|发布文章
| GET    |/api/v1/articles/:id|文章详情
| PUT    |/api/v1/articles/:id|修改文章
| DELETE |/api/v1/articles/:id|删除文章
| GET    |/api/v1/articles/:id/comments|查看文章评论
| POST   |/api/v1/articles/:id/comments|添加文章评论

## 架构
之前有一版是借鉴了著名Go开源项目 [Grafana](https://github.com/grafana/grafana) 的设计，使用了依赖注入机制，但是感觉过于复杂，不容易理解和使用，所以又改了。

本着简单易用易修改的原则，采用了**包全局变量**的方式初始化配置、日志、DB连接，项目目录如下：
```
├── config //配置
├── controller //控制器
├── log //日志
├── middleware //中间件
├── model //数据表模型
├── router //路由
├── service //服务
└── util //工具函数
```
在 ```main.go``` 里面依次**显示**初始化各个组件，清晰明了，简单易懂。

## 模块
### 1.配置
往往框架启动第一件事就是加载配置，Go的配置形式有很多种，比如最简单的JSON，但是JSON存在一个很大的问题就是无法注释，其次每增加一个配置就需要在结构体上面增加成员变量，也很麻烦。

我这里使用了ini作为配置文件格式，采用第三方库[gopkg.in/ini.v1](gopkg.in/ini.v1)来解析。

项目对这三方库的对象进行了简单包装，方便后续扩展，提供了一个Get()方法给外部调用，具体的用法可以参考库的文档。
```go
var cfg *App

func Get() *App {
    return cfg
}

type App struct {
    Env        string
    HttpPort   string
    LogFile    string
    LogConsole bool
    LogLevel   string
    
    *ini.File
}
```
另外，推荐不同环境采用不同的配置，比如app_dev.ini就是开发配置，可以通过-conf指定配置文件，默认会使用app.ini配置。

### 2.日志
日志这块最大的改变就是实现了全链路日志，方便问题排查，采用了第三方库zap log，在此基础上进行了封装。

```go
var zapLogger *zap.Logger

type Logger struct {
    context.Context
    *zap.Logger
}

// WithCtx 带请求上下文的Logger，可以记录一些额外信息，比如traceId
func WithCtx(ctx context.Context) *Logger {
    return &Logger{ctx, zapLogger}
}
```
还有一个就是实现了gorm的的sql日志，把sql日志也加上了traceId，为了实现链路日志，需要在记录日志的时候传入context对象，另外在控制器到模型层之间的调用也需要显示的传递context对象。

### 3.数据库
这块就是采用了第三方库gorm，具体用法可以参考其官方文档，这里面比较大的改动就是对日志的接口的实现，把sql日志记录下来了，详细可以看一下**gorm_logger.go**文件。

在配置里面支持多个数据库切换。
```go
var connPool = make(map[string]*gorm.DB)

// NewOrm 默认返回default数据库连接
func NewOrm(ctx context.Context, dbName ...string) *gorm.DB {
    conn := connPool["default"]
    if len(dbName) > 0 {
        if cn, ok := connPool[dbName[0]]; ok {
            conn = cn
        }
    }
    return conn.WithContext(ctx)
}
```

### 4.中间件
中间件这块添加了2个默认的，一个是增加链接日志用到的traceId，另一个记录一个访问日志，类似于Nginx的access_log，方便排查问题。

其余的话如果有需要可以自行增加。

### 5.MVC模型
由于接口并不涉及到视图文件，所以也不完全是MVC模型，不过我个人推荐3层结构：

controller ---> service ----> model

controller层的主要作用的对参数进行校验，项目采用了一个validation组件，通过结构体tag形式进行校验，校验通过后调用service层。

service层的主要作用是完成一些比较复杂的业务逻辑处理，不过很多增删改查接口逻辑本来就简单，这块也可以考虑去掉，直接调用model也不是不可。

model层的主要作用就是涉及数据库的操作。

这3层我都建议采用结构体成员函数的形式，然后使用全局变量初始化，这样用的时候就不用初始化了，直接调就行了，主要好处就是简单方便，当然这里注意一定不要在结构体里面存储私有数据，这样会出问题。
```go
type articleController struct {
    *Controller
    *service.ArticleService
}

var ArticleController = articleController{
    Controller:     BaseController,
    ArticleService: service.NewArticleService(),
}

// Create 添加文章
func (r articleController) Create(ctx *gin.Context) {
    var param model.CreateArticleCommand
    err := ctx.ShouldBindJSON(&param)
    if err != nil {
        r.Failed(ctx, ParamError, translate(err))
        return
    }
    if article, err := r.ArticleService.Create(ctx, &param); err != nil {
        r.Failed(ctx, Failed, err.Error())
    } else {
        r.Success(ctx, "添加文章成功", article)
    }
    return
```

## 命名规范
个人建议参考以下规范：

- 表名、文件夹名使用单数，比如```article、phone```
- 文件夹名全部小写，多个单词的话直接相连。如 ```eventhandler```
- 文件名小写下划线，如 ```article_controller.go```
- 变量名驼峰，如 ```var userId int```，```pageSize := 10```，首字母是否大写根据实际需要（是否要公开）
- 结构体名驼峰，如 ```type userService struct ```，首字母是否大写根据实际需要（是否要公开）
- 函数名驼峰，如果 ```func getUserById()```，首字母是否大写根据实际需要（是否要公开）

其它规范建议以Goland内置标准为准，一般情况下IDE都会提示，建议遵循。

## 使用
直接git clone本项目，然后根据自己的需求删除多余的控制器、模型等文件，如果有特殊需求直接修改代码即可，没有什么限制，Go的代码本来就很简单，我这里也只是给大家一个参考。