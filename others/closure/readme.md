## 闭包
### 简介
在支持函数是一等公民的语言中，一个函数的返回值是另一个函数，被返回的函数可以访问父函数内的变量，
当这个被返回的函数在外部执行时，就产生了闭包。

### go语言中的闭包
####标准库 net/http
``` go
func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error) {
 return func(*Request) (*url.URL, error) {
  return fixedURL, nil
 }
}
```
在返回的函数中，引用了父函数（ProxyURL）的参数 fixedURL，因此这是闭包。

