# htxp

一个用于go-zero框架的统一响应处理包，可以在goctl生成的API项目中使用。

## 安装

```bash
go get github.com/linktomarkdown/htxp
```

## 使用方法

在handler中导入此包：

```go
import "github.com/linktomarkdown/htxp"
```

然后可以使用以下函数处理HTTP响应：

```go
// 处理成功响应
htxp.Success(w, data)

// 处理错误响应
htxp.Error(w, err)

// 自定义响应
htxp.Response(w, data, err)
```

## 响应格式

成功响应：

```json
{
  "code": 200,
  "msg": "success",
  "data": {...}
}
```

错误响应：

```json
{
  "code": -1,
  "msg": "错误信息"
}
```

## 在goctl模板中使用

修改handler.tpl文件，将导入路径从本地response包改为此包：

```go
import (
  "net/http"
  
  "github.com/zeromicro/go-zero/rest/httpx"
  "github.com/linktomarkdown/htxp"
  {{.ImportPackages}}
)
```