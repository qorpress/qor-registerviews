# qor-registerviews

##### qor-registerviews
用于解决在 go mod 模式下使用 qor-admin 框架开发，相关组件视图注册失败的问题。

可以返回指定依赖包的指定文件夹路径 

###### 使用方法 

```go
package main

import (
"fmt"


registerviews "github.com/snowlyg/qor-registerviews"
	
)

func main() {
	
    path := registerviews.DetectViewsDir("github.com/snowlyg", "go-tenancy", "config") 
    fmt.Println(path)
}

```


详情见项目 [snowlyg/go-tenancy](https://github.com/snowlyg/go-tenancy)
