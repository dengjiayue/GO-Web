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

### 优化后性能测试

* 自搭建gee框架100万次压测

```
Server Software:        
Server Hostname:        localhost
Server Port:            8888

Document Path:          /
Document Length:        58 bytes

Concurrency Level:      100
Time taken for tests:   33.981 seconds
Complete requests:      1000000
Failed requests:        0
Total transferred:      198000000 bytes
HTML transferred:       58000000 bytes
Requests per second:    29427.98 [#/sec] (mean)
Time per request:       3.398 [ms] (mean)
Time per request:       0.034 [ms] (mean, across all concurrent requests)
Transfer rate:          5690.18 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   0.5      1       9
Processing:     0    2   3.8      2     369
Waiting:        0    2   3.7      1     368
Total:          0    3   3.8      3     370

Percentage of the requests served within a certain time (ms)
  50%      3
  66%      3
  75%      4
  80%      4
  90%      4
  95%      5
  98%      6
  99%      7
 100%    370 (longest request)
 ```

* gin框架100万次压测:

```
Server Software:        
Server Hostname:        localhost
Server Port:            8888

Document Path:          /
Document Length:        57 bytes

Concurrency Level:      100
Time taken for tests:   53.913 seconds
Complete requests:      1000000
Failed requests:        0
Total transferred:      180000000 bytes
HTML transferred:       57000000 bytes
Requests per second:    18548.48 [#/sec] (mean)
Time per request:       5.391 [ms] (mean)
Time per request:       0.054 [ms] (mean, across all concurrent requests)
Transfer rate:          3260.47 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.9      2      19
Processing:     0    4   2.5      3      58
Waiting:        0    3   2.3      2      56
Total:          0    5   2.6      5      60

Percentage of the requests served within a certain time (ms)
  50%      5
  66%      6
  75%      6
  80%      7
  90%      8
  95%     10
  98%     13
  99%     15
 100%     60 (longest request)
```

* gin框架:每秒请求数:18548.48，传输速率为1372.75 Kbytes/sec，平均响应时间为14.086毫秒
* 自搭建gee框架:每秒请求数为29427.98,传输速率为3260.47 Kbytes/sec,平均响应时间为 5.391 ms
* 根据结果数据，自搭建gee框架相对于gin框架每秒请求数高出约 58.76%，平均响应时间较短约 37.03%，传输速率高出约 74.44%。
