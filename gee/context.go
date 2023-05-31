package gee

import (
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type H map[string]interface{}

// 定义统一返回函数
func (c *Context) Returnfunc(code int, msg string, data interface{}) {
	c.JSON(code, H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

type Context struct {
	// origin objects
	Writer http.ResponseWriter //Writer: 一个 http.ResponseWriter 类型的对象，用于将响应内容写入 HTTP 响应体。
	Req    *http.Request       //Req: 一个 *http.Request 类型的指针，表示客户端发起的 HTTP 请求。
	// request info
	Path   string            //一个字符串，表示请求的 URL 路径。
	Method string            //一个字符串，表示请求的 HTTP 方法。get/post...
	Params map[string]string //一个字符串到字符串的映射，表示 URL 查询参数或表单数据。
	// response info
	StatusCode int //一个整数，表示 HTTP 响应状态码。
	// middleware
	Handlers []HandlerFunc //handlers: 一个 HandlerFunc 类型的切片，表示当前请求需要经过的中间件函数列表。
	Index    int           //index: 一个整数，表示当前执行到的中间件函数在 handlers 中的下标。可以理解为handers的下标
	Engine   *Engine
	Userid   int
}

// HandlerFunc 是一个函数类型，它接受一个 *Context 类型的参数，表示当前请求的上下文信息。
type HandlerFunc func(*Context)

// NewContext: 一个工厂函数，用于创建一个 Context 类型的对象，并将其初始化为一个请求的上下文信息。
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		Index:  -1,
	}
}

//Next: 用于将请求传递给下一个中间件函数。
//类似于洋葱模型的处理方式，每一层中间件函数处理完后，需要将请求传递给下一层中间件函数，直到最后一个中间件函数处理完请求并返回响应。

// Next() 函数的作用就是将请求传递给下一个中间件函数。它的实现比较简单，主要分为以下两步：

// 将 index 的值加一，表示将要执行下一个中间件函数。
// 循环执行中间件函数列表中从 index 开始的函数，直到执行到最后一个中间件函数或者其中一个函数调用了 c.Next() 方法，停止循环并返回
func (c *Context) Next() {
	c.Index++
	s := len(c.Handlers)
	for ; c.Index < s; c.Index++ {
		c.Handlers[c.Index](c)
	}
}

func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code) //输出状态码200/404...
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := jsoniter.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
		return
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// name为文件名
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.Engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
		fmt.Fprintf(c.Writer, "404 NOT FOUND: %s\n", c.Req.URL)
	}
}
func (c *Context) Fail(StatusCode int, err string) {
	http.Error(c.Writer, err, StatusCode)
}

func (c *Context) FailJson(StatusCode int, err string) {
	c.JSON(StatusCode, H{
		"msg": err,
	})
}

// 解析用户的json数据
func (c *Context) Getjson(v any) error {
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		return err
	}
	if err = jsoniter.Unmarshal(body, v); err != nil {
		return err
	}
	return nil
}
