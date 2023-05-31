package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Middlewares = append(group.Middlewares, middlewares...)
}

// 设置 HTML 模板中使用的函数，接收一个 template.FuncMap 类型的参数，将其存储到 engine.funcMap 中。
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 用于加载 HTML 模板文件，接收一个字符串类型的参数 pattern，表示要加载的模板文件的路径。
// 该方法首先通过 template.New("") 创建了一个空的模板，然后通过 engine.funcMap 将模板中使用的函数添加到模板中，
// 最后调用 template.ParseGlob(pattern) 加载指定路径下的所有模板文件，并通过 template.Must() 将加载的模板文件包装成一个不可变的 template.Template 类型，
// 并将其存储到 engine.htmlTemplates 中。这样，在程序运行过程中，就可以通过 engine.htmlTemplates 调用 HTML 模板了。
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 这行代码使用了Go语言标准库中的strings包的HasPrefix函数来判断req.URL.Path是否以group.Prefix为前缀。
	for _, group := range engine.Groups {
		if strings.HasPrefix(req.URL.Path, group.Prefix) {
			middlewares = append(middlewares, group.Middlewares...)
		}
	}
	c := NewContext(w, req)
	c.Handlers = middlewares
	c.Engine = engine
	engine.Router.Handle(c)
}

// user{login,name}
type RouterGroup struct {
	Prefix      string        //父元素的字符串
	Middlewares []HandlerFunc //中间件函数
	Parent      *RouterGroup  //父元素的节点
	Engine      *Engine       //所有的分组公用一个engine以实现所有的函数的接口
}
type Engine struct {
	*RouterGroup
	Router        *Rrouter           //
	Groups        []*RouterGroup     //
	htmlTemplates *template.Template //
	funcMap       template.FuncMap   //
}

// 初始化一个engine节点用于实现整个框架的功能
func New() *Engine {
	engine := &Engine{Router: NewRouter()}
	engine.RouterGroup = &RouterGroup{Engine: engine}
	engine.Groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 在当前路由分组下创建一个子分组，可以通过传入前缀字符串来指定子分组的前缀，该函数会返回一个新的 RouterGroup 对象，表示创建的子分组。在创建子分组时，需要指定父分组和引擎节点。
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.Engine
	newGroup := &RouterGroup{
		Prefix: group.Prefix + prefix,
		Parent: group,
		Engine: engine,
	}
	engine.Groups = append(engine.Groups, newGroup)
	return newGroup
}

func (group *RouterGroup) AddRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.Prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.Engine.Router.AddRoute(method, pattern, handler) //添加路由映射
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.AddRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.AddRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// createStaticHandler() 函数创建了一个处理静态文件请求的 HandlerFunc。
// 它接收两个参数：相对路径 relativePath 和文件系统 fs。该函数返回另一个 HandlerFunc，该函数负责检查请求的文件是否存在并且是否可以访问
// ，如果可以访问，就使用 http.FileServer 将文件传输给客户端。
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.Prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static() 函数是将静态文件服务器注册到路由中的函数。
// 它接收两个参数：相对路径 relativePath 和静态文件的根目录 root。它调用 createStaticHandler() 函数创建一个 HandlerFunc，并将其注册到路由器中，以便在请求 URL 匹配特定模式时调用。
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}
