package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// print stack trace for debug
func trace(message string) string {
	//runtime.Callers函数来获取当前goroutine的调用堆栈，并将其存储在一个长度为32的uintptr数组中，最多存储32个调用者的地址。
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) //该函数的返回值n表示实际存储在数组中的调用者数目。
	//runtime.Callers函数的第一个参数指定从哪个调用者开始跟踪，
	//通常传入3，以跳过当前函数trace、调用trace的函数log.Printf和Recovery函数。

	var str strings.Builder
	//错误信息、堆栈跟踪等信息写入该字符串中。
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] { //for循环遍历pcs[:n]数组，获取每个函数调用者的信息，并将其写入str中
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			//该代码中的recover函数用于恢复panic，并将其转换为字符串类型的message。然后，使用log.Printf函数打印错误信息和堆栈跟踪信息，并使用c.Fail函数向客户端发送一个HTTP错误响应。
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
				// fmt.Fprintf(c.Writer, "404 NOT FOUND: %s\n", c.Req.URL)
			}
		}()
		c.Next()
	}
}
