package gee

import (
	"net/http"
	"strings"
)

type Rrouter struct {
	Roots    map[string]*Node
	Handlers map[string]HandlerFunc
}

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']

func NewRouter() *Rrouter {
	return &Rrouter{
		Roots:    make(map[string]*Node),
		Handlers: make(map[string]HandlerFunc),
	}
}

// user/login/:id==>{user,login,:id(123)}
func ParsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
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

func (r *Rrouter) AddRoute(method string, pattern string, handler HandlerFunc) {
	parts := ParsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.Roots[method]
	if !ok {
		r.Roots[method] = &Node{}
	}
	r.Roots[method].Insert(pattern, parts, 0)
	r.Handlers[key] = handler
}

func (r *Rrouter) GetRoute(method string, path string) (*Node, map[string]string) {
	searchParts := ParsePattern(path)
	params := make(map[string]string)
	root, ok := r.Roots[method]

	if !ok {
		return nil, nil
	}

	n := root.Search(searchParts, 0)

	if n != nil {
		parts := ParsePattern(n.Pattern)
		//提取动态路由 的路由参数在map中
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *Rrouter) Handle(c *Context) {
	n, params := r.GetRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + n.Pattern
		c.Params = params
		//添加路由函数
		c.Handlers = append(c.Handlers, r.Handlers[key])
	} else {
		c.Handlers = append(c.Handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND:method:%s path , %s\n", c.Method, c.Path)
		})
	}
	c.Next()
}
