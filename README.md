# Gen Web - 基于Gin框架封装的脚手架结构，便于快速开发API

## 介绍

主要使用以下开源组件：

* [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [gorm.io/gorm](https://gorm.io/gorm)
* [go.uber.org/zap](https://go.uber.org/zap)
* [github.com/go-playground/validator/v10](https://github.com/go-playground/validator/v10)
* [github.com/facebookgo/inject](https://github.com/facebookgo/inject)
* [github.com/go-redis/redis](https://github.com/go-redis/redis)

目录结构清晰明了，项目包含了一个用户注册、登录、文章增删改查等功能的 Restful API 应用，仅供参考！

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

项目采用了依赖注入的方式贯穿全局，我们可以把DB、缓存、HTTP API等功能看作是项目的一个服务，通过facebook开源的inject库，我们在启动项目把这些```Service```注入进去，解决各自之间的依赖关系。

```go
type ArticleService struct {
SQLStore *SQLService         `inject:""`
Cache    *cache.CacheService `inject:""`
}

func init() {
registry.RegisterService(&ArticleService{})
}

func (r ArticleService) Init() error {
return nil
}
```

既灵活，也不影响性能，因为虽然依赖注入使用了反射，但是我们只在程序启动之前做这件事，而且只需要进行一次。

## 启动流程

```main```文件是程序的入口，主要功能是解析命令行参数，只有一个参数，那就是配置文件，默认配置文件是当前目录下的```app.ini```

紧接着，创建一个```Server```实例：

```go
// Server is responsible for managing the lifecycle of services.
type Server struct {
context          context.Context
shutdownFn       context.CancelFunc
childRoutines    *errgroup.Group
log              *zap.Logger
cfg              *config.Cfg // 项目配置
shutdownOnce     sync.Once
shutdownFinished chan struct{}
isInitialized    bool
mtx              sync.Mutex
serviceRegistry  serviceRegistry // 注册的服务
}
```

这个Server实例是管理所有服务的中心，其主要工作就是加载配置文件，然后根据配置文件初始化日志配置，日志库采用zap log，主要文件在```zap/zap_logger.go```里面

然后还有一个最重要是就是初始化所有注册过服务，执行其```Init```方法做一些初始化工作，最后执行后台服务。

如果一个服务实现了```Run```方法，就是一个后台服务，会在项目启动时候运行，结束时候优雅关闭，其中最好的例子就是```HTTPServer```，我们可以把API服务认为是一个后台服务，在整个项目启动的时候就会运行。

```go
type HTTPServer struct {
log     *zap.Logger
gin     *gin.Engine
context context.Context

Cfg            *config.Cfg             `inject:""`
ArticleService *article.ArticleService `inject:""`
UserService    *user.UserService       `inject:""`
}
```

HTTPServer的代码在```api/http_server.go```文件里面，其主要作用就是初始化一些服务配置，然后启动HTTP服务，使用了Gin框架。

## 代码介绍

在```services```文件夹下包含了一些服务的代码文件。

项目整体是一个3层架构，即控制器层、Service层、模型层。

个人理解，控制器层主要做一些接口参数校验等工作，模型层主要是数据操作，Service层才是主要的业务逻辑。

数据库相关配置在```models/db.go```里面，也是一个服务，主要作用是根据配置文件初始化数据库连接，支持多数据库切换、支持Sql日志记录。

```go
type SQLService struct {
Cfg *config.Cfg `inject:""`

conns map[string]*gorm.DB
log   *zap.Logger
}

func DB(dbName ...string) *gorm.DB {
if len(dbName) > 0 {
if conn, ok := sqlStore.conns[dbName[0]]; ok {
return conn
}
}
return db
}
```

项目使用了Gorm（2.0版本），具体详细用法可以参考官方文档。

路由文件位于```api/api.go```，可以多层嵌套，中间件在```middleware```文件夹。

```config/config.go```是配置文件的一些加载逻辑，可以根据自己需求适当的修改优化。

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

## 使用

建议直接clone本项目，然后删除多余的控制器、模型等文件，根据自己需求调整即可。

```go
debug   registry/di.go:42       Service [SqlService] init success       {"time": "2021-06-23 23:59:15"}
debug   registry/di.go:42       Service [HTTPServer] init success       {"time": "2021-06-23 23:59:15"}
debug   registry/di.go:42       Service [UserService] init success      {"time": "2021-06-23 23:59:15"}
debug   registry/di.go:42       Service [CacheService] init success     {"time": "2021-06-23 23:59:15"}
debug   registry/di.go:42       Service [ArticleService] init success   {"time": "2021-06-23 23:59:15"}
debug   gen/main.go:39  Waiting on services...  {"time": "2021-06-23 23:59:15", "module": "server"}
debug   server/server.go:126    server was started successfully {"time": "2021-06-23 23:59:15", "module": "http_server"}
info    sync/once.go:66 Shutdown started        {"time": "2021-06-23 23:59:15", "module": "server"}
debug   server/server.go:126    server was shutdown gracefully  {"time": "2021-06-23 23:59:15", "module": "http_server"}
debug   errgroup/errgroup.go:57 Stopped HTTPServer      {"time": "2021-06-23 23:59:15", "module": "server"}
``` 

最后的最后，本项目参考借鉴了著名Go开源项目 [Grafana](https://github.com/grafana/grafana) 的设计和架构，这个项目的后端是全部采用Go开发，东西也很多，代码很不错。