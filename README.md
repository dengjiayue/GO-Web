# Go-Web
* 描述:自搭建 goweb 框架并优化


### 思想
web框架服务主要围绕着请求与响应来展开的

搭建一个web框架的核心思想

1 便捷添加响应路径与响应函数(base)

2 能够接收多种数据类型传入(上下文context)

3 构建多种数据类型的响应(上下文context)

4 使用前缀树储存搜索路径,实现动态路由的功能(前缀树)

5 能够进行分组与中间件的插入(分组,中间件)

### 优化

1 替换原生的json数据解析函数

使用JSON-iterator对json数据进行解析(编码与解码)

2 构建线程池限制同时用户访问的数量

3 搭建统一数据返回的函数,使用
```json
{
code:状态码
msg:提示信息
data:返回数据
}
```
的格式进行返回,便于前端使用

4 搭建传入json数据的解析,解析任意形式的json的数据,便于数据的读入
(个人比较喜欢使用json格式的数据,将请求与响应的数据统一为json类型,前端后端都比较分别使用)