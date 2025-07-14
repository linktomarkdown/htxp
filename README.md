# HTXP - Go工具库

HTXP 是一个 Go 语言工具库，提供了常用的功能封装，包括字符串生成、文件操作、Minio对象存储、RabbitMQ消息队列、MDM管理等。

## 特性

- 🚀 **简洁的API设计** - 使用包级别函数，无需实例化
- 🔧 **丰富的工具函数** - 字符串生成、加密、文件操作等
- ☁️ **对象存储支持** - Minio集成
- 📨 **消息队列支持** - RabbitMQ集成
- 👥 **用户管理** - Minio MDM用户管理
- 📝 **日志系统** - Logrus集成
- 🔄 **HTTP响应** - 统一的HTTP响应格式

## 安装

```bash
go get github.com/linktomarkdown/htxp
```

## 快速开始

### 基础工具函数

```go
package main

import (
    "fmt"
    "github.com/linktomarkdown/htxp"
)

func main() {
    // 初始化日志
    htxp.InitLogrus()
    
    // 生成订单号
    orderNo := htxp.GenerateOrderNo()
    fmt.Println("订单号:", orderNo)
    
    // 生成微信支付订单ID
    orderID := htxp.GenerateOrderID("wechat")
    fmt.Println("订单ID:", orderID)
    
    // 生成随机名称
    name := htxp.GenerateName(10)
    fmt.Println("随机名称:", name)
    
    // MD5加密
    hash := htxp.Md5V("password123")
    fmt.Println("MD5哈希:", hash)
    
    // 生成随机密码
    password := htxp.GenerateRandomPassword(12, true, true, true)
    fmt.Println("随机密码:", password)
}
```

### Minio对象存储

```go
package main

import (
    "log"
    "github.com/linktomarkdown/htxp"
)

func main() {
    // 初始化Minio客户端
    err := htxp.InitMinio("localhost:9000", "minioadmin", "minioadmin", false)
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建bucket
    err = htxp.CreateBucket("mybucket")
    if err != nil {
        log.Fatal(err)
    }
    
    // 检查bucket是否存在
    exists, err := htxp.BucketHasExists("mybucket")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Bucket exists: %v", exists)
}
```

### RabbitMQ消息队列

```go
package main

import (
    "log"
    "github.com/linktomarkdown/htxp"
)

func main() {
    // 初始化RabbitMQ连接
    rabbitConf := htxp.RabbitConf{
        Username: "guest",
        Password: "guest", 
        Host:     "localhost",
        Port:     5672,
        VHost:    "/",
    }
    err := htxp.InitRabbitMQ(rabbitConf)
    if err != nil {
        log.Fatal(err)
    }
    
    // 初始化交换机和队列
    err = htxp.InitRabbitMQExchange("test_exchange", "direct")
    if err != nil {
        log.Fatal(err)
    }
    
    // 发送消息
    err = htxp.SendMessage("test_exchange", "test_key", []byte("Hello RabbitMQ"), "text/plain")
    if err != nil {
        log.Fatal(err)
    }
    
    // 关闭连接
    defer htxp.CloseRabbitMQ()
}
```

### MDM用户管理

```go
package main

import (
    "log"
    "github.com/linktomarkdown/htxp"
)

func main() {
    // 初始化MDM客户端
    err := htxp.InitMDM("localhost:9000", "minioadmin", "minioadmin", false)
    if err != nil {
        log.Fatal(err)
    }
    
    // 创建用户
    err = htxp.CreateUser("newuser", "newpassword")
    if err != nil {
        log.Fatal(err)
    }
    
    // 设置用户策略
    err = htxp.SetUserPolicy("newuser", "readwrite")
    if err != nil {
        log.Fatal(err)
    }
    
    // 获取用户信息
    userInfo, err := htxp.UserInfo("newuser")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("用户信息: %+v", userInfo)
}
```

### HTTP响应

```go
package main

import (
    "errors"
    "net/http"
    "github.com/linktomarkdown/htxp"
)

func main() {
    http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]interface{}{
            "name": "张三",
            "age":  30,
        }
        htxp.Success(w, data)
    })

    http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
        err := errors.New("发生了一个错误")
        htxp.Error(w, err)
    })

    http.ListenAndServe(":8080", nil)
}
```

## API文档

### 字符串生成

- `GenerateOrderNo()` - 生成订单号
- `GenerateOrderID(paymentType string)` - 生成支付订单ID
- `GenerateName(n int)` - 生成随机名称
- `GenerateRandomPassword(length int, useLetters, useSpecial, useNum bool)` - 生成随机密码
- `GenerateRandomString(length int)` - 生成随机字符串
- `GenerateRandomNumber(length int)` - 生成随机数字
- `GenerateRandomUUID()` - 生成UUID

### 加密工具

- `Md5V(str string)` - MD5加密
- `GenerateKey(length int)` - 生成密钥

### 类型转换

- `StringToInt(s string)` - 字符串转整数
- `StringToFloat64(s string)` - 字符串转浮点数
- `ConvertUidToUint64(uid string)` - UID转uint64

### 数组操作

- `InArray(needle string, haystack []string)` - 判断元素是否在数组中
- `ContainsRoles(needle string, haystack []string)` - 判断是否包含角色

### 文件操作

- `CopyFile(src, dst string)` - 复制文件
- `CopyDir(src, dst string)` - 复制目录

### 日志系统

- `InitLogrus()` - 初始化日志系统

### 环境变量

- `GetEnvInfo(env string)` - 获取环境变量

### 异常处理

- `TryCatch(f func(), handler func(interface{}))` - 异常捕获

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License