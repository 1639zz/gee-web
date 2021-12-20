package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots['GET'] roots['POST']
// handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc)}
}

//分割字符串（只有*适合）
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	//遍历循环判断子节点是否有* 有则停下
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//获取分割后的子节点
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	//得到路由信息中的method方法
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	//插入路由信息
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

//获取路由信息
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//获取url路径
	searchParts := parsePattern(path)
	//得到params
	params := make(map[string]string)
	//得到root值
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	//搜索根节点
	n := root.search(searchParts, 0)
	//当搜索后的根节点不为空时，对子节点分割后进行判断
	if n != nil {
		parts := parsePattern(n.pattern)
		//循环遍历判断
		for index, part := range parts {
			//如果第一个为':',则根据下标搜索子节点
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			//如果为'*' 或者子节点的长度>1
			if part[0] == '*' && len(part) > 1 {
				//根据下标进行搜索，并通过/进行拼接
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

//func (r *router) handle(c *Context) {
//	key := c.Method + "-" + c.Path
//	if handler, ok := r.handlers[key]; ok {
//		handler(c)
//	} else {
//		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
//	}
//}
//根据不同的请求判断适用于哪个中间件，得到中间件list列表
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	//如果不为空
	if n != nil {
		c.Params = params
		//得到key参数
		key := c.Method + "-" + n.pattern
		//添加到handlers数组
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
