package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP

//添加RouterGroup接口
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // store all groups
}

//RouterGroup struct
type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	//engine继承RouterGroup所有方法
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//创建中间件数组
	var middlewares []HandlerFunc
	//循环遍历所有group
	for _, group := range engine.groups {
		//如果字符串hasprefix
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			//添加中间件信息
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

//解析请求地址
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//得到地址，并放在group中
	absolutePath := path.Join(group.prefix, relativePath)
	//保存fileServer
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		//检查文件是否创建、是否有权限进行操作
		if _, err := fs.Open(file); err != nil {
			//返回404错误
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

//静态文件处理(暴露给用户)
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}
